package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"

	"image-sync/dockerhub"

	registryV1 "github.com/oscarzhao/docker-reg-client/registry" // v1
	registryV2 "github.com/oscarzhao/docker-registry-client/registry"
)

type Client struct {
	isHub       bool
	proto       string
	registry    string
	version     string
	repoName    string
	repoTag     string
	RegClient   *registryV1.Client
	RegClientV2 *registryV2.Registry
	HubClient   *dockerhub.DockerHubClient
}

// NewClient creates a new registry client, default returns a docker hub client
func NewClient(proto, registry, version, username, password string) (*Client, error) {
	if registry == "" || version == "" || proto == "" || registry == "index.docker.io" {
		return &Client{
			isHub:     true,
			proto:     "https",
			registry:  "index.docker.io",
			version:   "v2",
			HubClient: &dockerhub.DockerHubClient{},
		}, nil
	}
	switch version {
	case "v1":
		srcClient, err := registryV1.NewClient(proto, registry)
		if err != nil {
			return nil, err
		}
		return &Client{
			isHub:     false,
			proto:     proto,
			registry:  registry,
			version:   version,
			RegClient: srcClient,
		}, nil
	case "v2":
		srcClient, err := registryV2.New(fmt.Sprintf("%s://%s/", proto, registry), username, password)
		if err != nil {
			return nil, err
		}
		return &Client{
			isHub:       false,
			proto:       proto,
			registry:    registry,
			version:     version,
			RegClientV2: srcClient,
		}, nil
	}

	return nil, errors.New("invalid client config")
}

// IsHub returns true if c is a docker hub client, false otherwise
func (c *Client) IsHub() bool {
	return c.isHub
}

// ListRepositories list all repos according to a keyword
func (c *Client) ListRepositories(pattern string) ([]string, error) {
	if c.isHub {
		images, err := c.HubClient.SearchReposByUser(pattern)
		if err != nil {
			return nil, err
		}
		var res []string
		for _, image := range images {
			res = append(res, image.RepoName)
		}
		return res, nil
	}

	switch c.version {
	case "v1":
		return c.ListRepositoriesV1(pattern)
	case "v2":
		return nil, errors.New("registry v2 does not support search repo")
	default:
		return nil, errors.New("invalid registry version")
	}
}

// ListRepositoriesV1 lists all repos according to the pattern
func (c Client) ListRepositoriesV1(pattern string) ([]string, error) {
	repoList := make([]string, 0, 64)

	searchResults, err := c.RegClient.Search.Query(pattern, 0, 100)
	if err != nil {
		return nil, err
	}

	for _, res := range searchResults.Results {
		if pattern != "library" {
			if strings.HasPrefix(res.Name, pattern) {
				repoList = append(repoList, res.Name)
			}
		} else {
			arr := strings.Split(res.Name, "/")
			if len(arr) >= 2 {
				continue
			}
			repoList = append(repoList, "library/"+arr[0])
		}
	}

	pageNumber := searchResults.Page
	glog.V(7).Infof("ListRepositories, pages:%d, page_size:%d\n", pageNumber, searchResults.PageSize)

	for i := 1; i < pageNumber; i++ {
		tempResult, err := c.RegClient.Search.Query(pattern, i, 100)
		if err != nil {
			glog.Errorf("Search repo failed, pattern: %s, page:%d, err:%s\n", pattern, i, err)
			return repoList, err
		}

		for _, res := range tempResult.Results {
			if pattern != "library" {
				if strings.HasPrefix(res.Name, pattern) {
					repoList = append(repoList, res.Name)
				}
			} else {
				arr := strings.Split(res.Name, "/")
				if len(arr) >= 2 {
					continue
				}
				repoList = append(repoList, "library/"+arr[0])
			}
		}

	}
	return repoList, nil
}

// ListTags lists all tags of a repo in certain registry server
func (c Client) ListTags(repo string) ([]string, error) {
	if c.isHub {
		tags, err := c.HubClient.QueryImageTags(repo)
		if err != nil {
			return nil, err
		}
		var res []string
		for _, t := range tags {
			res = append(res, t.Name)
		}
		return res, nil
	}

	if c.version == "v1" {
		auth, err := c.RegClient.Hub.GetReadToken(repo)
		if err != nil {
			glog.Errorf("GetReadToken failed:%s\n", err)
			return nil, err
		}

		tagMap, err := c.RegClient.Repository.ListTags(repo, auth)
		if err != nil {
			glog.Errorf("ListTags failed, error info:%s\n", err)
			return nil, err
		}
		glog.V(6).Infof("ListTags v1 succeeds, repo:%s, results: %v\n", repo, tagMap)
		tags := make([]string, 0, 64)
		for key := range tagMap {
			tags = append(tags, key)
		}
		return tags, nil
	}

	// use v2 version api
	tags, err := c.RegClientV2.Tags(repo)
	if err != nil {
		glog.Errorf("ListTags failed, error info: %s\n", err)
		return nil, err
	}
	glog.V(6).Infof("ListTags v2 succeeds, repo: %s, results: %v\n", repo, tags)
	return tags, nil
}
