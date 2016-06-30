package registry

import (
	"strings"
	"testing"
)

type RepositoryPair struct {
	SrcRepoUrl string
	DstRepoUrl string
}

type imageInfo struct {
	proto    string
	registry string
	version  string
	repoName string
	tag      string
	username string
	password string
}

func TestListRepositories(t *testing.T) {
	shouldSuccesses := []imageInfo{
		{"", "", "", "oscarzhao", "latest", "", ""}, // docker hub search repo is supported
		// howerver, firewall make this fails
		// {"https", "gcr.io", "v1", "google_containers/pause", "latest", "", ""}, // registry v1 search repo is supported
	}

	for _, tc := range shouldSuccesses {
		srcClient, err := NewClient(tc.proto, tc.registry, tc.version, tc.username, tc.password)
		if err != nil {
			t.Errorf("ListRepositories, failed to create client, error info:%s\n", err)
			continue
		}
		repos, err := srcClient.ListRepositories(tc.repoName)
		if err != nil {
			t.Errorf("search repo failed, tc config:%#v, error:%s\n", tc, err)
		} else {
			t.Logf("search success: %d\n%s\n", len(repos), strings.Join(repos, ", "))
		}
	}

	shouldFails := []imageInfo{
		{"https", "index.tenxcloud.com", "v1", "docker_library/notfound", "latest", "", ""}, // registry version(v1) not match the actual version(v2)
		{"https", "index.tenxcloud.com", "v2", "google_containers/pause", "latest", "", ""}, // registry v2 search is not supported
	}
	for _, tc := range shouldFails {
		srcClient, err := NewClient(tc.proto, tc.registry, tc.version, tc.username, tc.password)
		if err != nil {
			t.Errorf("failed to create client, config:%#v error info:%s\n", tc, err)
			continue
		}
		repos, err := srcClient.ListRepositories(tc.repoName)
		if err == nil {
			t.Errorf("Should report error, but returned result:%#v\n", repos)
		}
	}

}

func TestListTags(t *testing.T) {
	shouldSuccesses := []imageInfo{
		{"https", "", "v2", "alpine", "latest", "", ""},
		// {"https", "index.tenxcloud.com", "v2", "google_containers/pause", "latest", "", ""},
		// {"https", "index.docker.io", "v2", "google/golang", "latest", "", ""},
		// {"https", "index.docker.io", "v1", "google/notfound", "latest", "", ""},
	}

	for _, tc := range shouldSuccesses {
		srcClient, err := NewClient(tc.proto, tc.registry, tc.version, tc.username, tc.password)
		if err != nil {
			t.Errorf("list repos, failed to create client, error info:%s\n", err)
			continue
		}
		repos, err := srcClient.ListTags(tc.repoName)
		if err != nil {
			t.Errorf("list repos failed, error:%s\n", err)
		} else {
			t.Logf("list repos success, repo:%s, tag count:%d\n%s\n", tc.repoName, len(repos), strings.Join(repos, ", "))
		}
	}
}
