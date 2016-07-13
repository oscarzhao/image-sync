package main

import (
	// "errors"
	"flag"
	"strings"

	"github.com/golang/glog"

	"github.com/oscarzhao/image-sync/dockerexec"
	"github.com/oscarzhao/image-sync/registry"
)

var (
	srcClient *registry.Client
	dstClient *registry.Client

	// registry configs
	srcRegistry        string
	srcRegistryVersion string
	srcRepoPassword    string
	dstRegistry        string
	dstRegistryVersion string
	dstRepoPassword    string

	srcRepoOwner string
	dstRepoOwner string
)

type Image struct {
	registry string
	repo     string
	tag      string
}

func (i Image) String() string {
	if i.registry == "" {
		return i.repo + ":" + i.tag
	}
	return i.registry + "/" + i.repo + ":" + i.tag
}

func init() {
	flag.Set("alsologtostderr", "true")
	flag.StringVar(&srcRegistry, "src-registry", "", "use docker hub as default, alternatives: gcr.io")
	flag.StringVar(&srcRegistryVersion, "src-registry-version", "v2", "the registry api version (v1 or v2)")
	flag.StringVar(&srcRepoPassword, "src-repo-password", "xxx", "repo password, currently not useful")

	flag.StringVar(&dstRegistry, "dst-registry", "index.tenxcloud.com", "the registry to synchronize to")
	flag.StringVar(&dstRegistryVersion, "dst-registry-version", "v2", "the registry api version (often v2)")
	flag.StringVar(&dstRepoPassword, "dst-repo-password", "xxx", "repo password, use to list repos at dst registry")

	flag.StringVar(&srcRepoOwner, "repo-owner", "", "repo owner, the user images are under for the source registry")
	flag.StringVar(&dstRepoOwner, "dst-repo-owner", "docker_library", "repo owner, the user images are under for the target registry")
	flag.Parse()

	if srcRepoOwner == "" || srcRepoOwner == "library" {
		srcRepoOwner = "library" // empty is library
		dstRepoOwner = "docker_library"
	}

	srcClient, _ = registry.NewClient("https", srcRegistry, srcRegistryVersion, srcRepoOwner, srcRepoPassword)
	dstClient, _ = registry.NewClient("https", dstRegistry, dstRegistryVersion, dstRepoOwner, dstRepoPassword)
}

func main() {
	srcRepo2Tags := make(map[string][]string)
	listTagFailedRepos := make([]string, 0, 4)

	repoList, err := srcClient.ListRepositories(srcRepoOwner)
	if err != nil {
		glog.Errorf("list repos (%s) failed, error: %s\n", srcRepoOwner, err)
		return
	}

	glog.V(4).Infof("repos got: %s\n", strings.Join(repoList, "\n"))

	// fetch all tags of all repos under srcRepoOwner
	for _, repoName := range repoList {
		tags, err := srcClient.ListTags(repoName)
		if err != nil {
			listTagFailedRepos = append(listTagFailedRepos, repoName)
			glog.Errorf("list tag of repo (%s/%s) fails, error:%s\n", srcRegistry, repoName, err)
			continue
		}
		srcRepo2Tags[repoName] = tags
	}

	glog.V(4).Infof("images found in source registry: %#v\n", srcRepo2Tags)
	images2pull := listImagesToPull(srcRepo2Tags)
	imagePulled := pullImages(images2pull)
	images2push := makeTag(imagePulled, dstRegistry)
	imagePushed := pushImages(images2push)

	for pushed := range imagePushed {
		glog.V(2).Infof("image %s pushed\n", pushed)
	}
	if len(listTagFailedRepos) > 0 {
		glog.Errorf("the following repos, list tag operation fails:\n%s\n", strings.Join(listTagFailedRepos, ", "))
	}
}

func listImagesToPull(repo2tags map[string][]string) <-chan Image {
	images2pull := make(chan Image)
	go func() {
		for repo, tags := range repo2tags {
			for _, tag := range tags {
				images2pull <- Image{registry: srcRegistry, repo: repo, tag: tag}
			}
		}
		close(images2pull)
	}()
	return images2pull
}

func pullImages(images <-chan Image) <-chan Image {
	success := make(chan Image)
	go func() {
		for image := range images {
			if _, stderr, err := dockerexec.PullImage(image.registry, image.repo, image.tag); err != nil {
				glog.Errorf("dockerexec.PullImage (%v) failed, stderr:%s, err:%s\n", image, stderr, err)
			} else {
				success <- image
			}
		}
		close(success)
	}()
	return success
}

func pushImages(images <-chan Image) <-chan Image {
	success := make(chan Image)
	go func() {
		for image := range images {
			if _, stderr, err := dockerexec.PushImage(image.registry, image.repo, image.tag); err != nil {
				glog.Errorf("dockerexec.PushImage %v failed, stderr:%s, err:%s, mark and delete it\n", image, stderr, err)
				go func(registry, repo, tag string) {
					if _, stderr, err := dockerexec.DeleteImage(registry, repo, tag); err != nil {
						glog.Errorf("delete image %s/%s:%s fails, stderror:%s, error:%s\n", registry, repo, tag, stderr, err)
					}
				}(image.registry, image.repo, image.tag)
			} else {
				success <- image
			}
		}
		close(success)
	}()
	return success
}

func makeTag(images <-chan Image, dstRegistry string) <-chan Image {
	success := make(chan Image)
	go func() {
		for image := range images {
			// check if create tag success
			dstRepo := image.repo
			if image.registry == "" {
				dstRepo = dstRepoOwner + "/" + dstRepo
			}
			dstImg := Image{dstRegistry, dstRepo, image.tag}
			if _, stderr, err := dockerexec.MakeTag(image.String(), dstImg.String()); err == nil {
				success <- dstImg
			} else {
				glog.Errorf("create tag from %s to %s fails, stderr:%s, error:%s\n", image, dstImg, stderr, err)
			}
			// delete old one
			if _, stderr, err := dockerexec.DeleteImage(image.registry, image.repo, image.tag); err != nil {
				glog.Errorf("delete image %s fails, stderror:%s, error:%s\n", image, stderr, err)
			}
		}
		close(success)
	}()

	return success
}
