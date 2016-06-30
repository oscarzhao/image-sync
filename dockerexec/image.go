package dockerexec

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

// PullImage pulls image from a registry server
func PullImage(registry, repo, tag string) (stdout, stderr string, err error) {
	var stdoutB, stderrB bytes.Buffer
	var image string
	registry = strings.Trim(registry, "/")
	if registry == "" {
		image = fmt.Sprintf("%s:%s", repo, tag)
	} else {
		image = fmt.Sprintf("%s/%s:%s", registry, repo, tag)
	}
	cmd := exec.Command("/usr/bin/docker", "pull", image)
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err = cmd.Run()
	stdout = stdoutB.String()
	stderr = stderrB.String()
	return
}

// PushImage pushes image to a registry server
func PushImage(registry, repo, tag string) (stdout, stderr string, err error) {
	var stdoutB, stderrB bytes.Buffer
	var image string
	registry = strings.Trim(registry, "/")
	if registry == "" {
		image = fmt.Sprintf("%s:%s", repo, tag)
	} else {
		image = fmt.Sprintf("%s/%s:%s", registry, repo, tag)
	}
	cmd := exec.Command("/usr/bin/docker", "push", image)
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err = cmd.Run()
	stdout = stdoutB.String()
	stderr = stderrB.String()
	return
}

// DeleteImage deletes image from a registry server
func DeleteImage(registry, repo, tag string) (stdout, stderr string, err error) {
	var stdoutB, stderrB bytes.Buffer
	var image string
	registry = strings.Trim(registry, "/")
	if registry == "" {
		image = fmt.Sprintf("%s:%s", repo, tag)
	} else {
		image = fmt.Sprintf("%s/%s:%s", registry, repo, tag)
	}
	cmd := exec.Command("/usr/bin/docker", "rmi", image)
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err = cmd.Run()
	stdout = stdoutB.String()
	stderr = stderrB.String()
	return
}

// MakeTag creates a new tag from an existing image
func MakeTag(from, to string) (stdout, stderr string, err error) {
	var stdoutB, stderrB bytes.Buffer
	cmd := exec.Command("/usr/bin/docker", "tag", from, to)
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err = cmd.Run()
	stdout = stdoutB.String()
	stderr = stderrB.String()
	return
}

// ListLocalImageAndTags lists all images and tags in local disk
func ListLocalImageAndTags() (map[string][]string, error) {
	var stdoutB, stderrB bytes.Buffer
	cmd := exec.Command("/usr/bin/docker", "images")
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err := cmd.Run()
	stdout := stdoutB.String()
	stderr := stderrB.String()
	if err != nil {
		glog.Errorf("ListLocalImageAndTags failed, stderr:%s, err:%s\n", stderr, err)
		return nil, errors.New(stderr)
	}

	pattern := `(\S+)([ \t\r\f]+)(\S+).+`
	re := regexp.MustCompile(pattern)

	imageArr := strings.Split(string(stdout), "\n")
	imageArr = imageArr[1:]
	image2tags := make(map[string][]string)
	for _, imgStr := range imageArr {
		results := re.FindStringSubmatch(imgStr)
		if len(results) < 4 {
			fmt.Printf("ListLocalImageAndTags, invalid row: %s\n", imgStr)
			continue
		}
		repo := results[1]
		tag := results[3]
		if _, ok := image2tags[repo]; ok {
			image2tags[repo] = append(image2tags[repo], tag)
		} else {
			image2tags[repo] = []string{tag}
		}
	}
	return image2tags, nil
}
