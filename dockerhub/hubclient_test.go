package dockerhub

import (
	"testing"
)

var (
	c = &DockerHubClient{}
)

func TestSearchReposByUser(t *testing.T) {
	user := "oscarzhao"
	images, err := c.SearchReposByUser(user)

	if err != nil {
		t.Errorf("list %s's images fails, should succeeed, error:%s\n", user, err)
		return
	}

	t.Logf("list %s's images success\n%#v\n", user, images)
}

func TestListTags(t *testing.T) {
	repo := "ubuntu"
	tags, err := c.QueryImageTags(repo)

	if err != nil {
		t.Errorf("list %s's images fails, should succeeed, error:%s\n", repo, err)
		return
	}

	t.Logf("list %s's images success, count: %v\n", repo, len(tags))
}
