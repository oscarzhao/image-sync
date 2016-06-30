package dockerexec

import (
	"strings"
	"testing"
)

func TestPullImage(t *testing.T) {
	shoudSuccess := []struct {
		registry string
		repo     string
		tag      string
	}{
		{"", "hello-world", "latest"},
	}
	shouldFailure := []struct {
		registry string
		repo     string
		tag      string
	}{
		{"", "not-found", "not-found"},
	}

	// test success
	for _, img := range shoudSuccess {
		if _, stderr, err := PullImage(img.registry, img.repo, img.tag); err != nil {
			t.Errorf("pull image %#v should succeed, but failed. stderr:%s, error:%s\n", img, stderr, err)
		}
	}

	// test failed
	for _, img := range shouldFailure {
		if _, _, err := PullImage(img.registry, img.repo, img.tag); err == nil {
			t.Errorf("pull image %#v should fail, but success\n", img)
		}
	}

	// delete image pulled
	for _, img := range shoudSuccess {
		if _, stderr, err := DeleteImage(img.registry, img.repo, img.tag); err != nil {
			t.Errorf("delete image %#v should succeed, but failed. stderr:%s, error:%s\n", img, stderr, err)
		}
	}
}

func TestMakeTag(t *testing.T) {
	shoudSuccess := []struct {
		from string
		to   string
	}{
		{
			from: "hello-world:latest",
			to:   "index.tenxcloud.com/docker_library/hello-world:new-tag",
		},
	}
	for _, tags := range shoudSuccess {
		arr := strings.Split(tags.from, ":")
		_, stderr, err := PullImage("", arr[0], arr[1])
		if err != nil {
			t.Errorf("pull image %s fails, stderr: %s, error:%s\n", tags.from, stderr, err)
			continue
		}
		_, stderr, err = MakeTag(tags.from, tags.to)
		if err != nil {
			t.Errorf("make tag should succeed, but fails, stderr:%s, err:%s\n", stderr, err)
			continue
		}
		// delete image
		_, stderr, err = DeleteImage("", arr[0], arr[1])
		if err != nil {
			t.Errorf("delete tag %s should succeed, but fails, stderr:%s, error:%s\n", tags.to, stderr, err)
		}
		// delete tag
		arr = strings.Split(tags.to, ":")
		_, stderr, err = DeleteImage("", arr[0], arr[1])
		if err != nil {
			t.Errorf("delete tag %s should succeed, but fails, stderr:%s, error:%s\n", tags.to, stderr, err)
		}
	}
}
