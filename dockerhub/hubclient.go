package dockerhub

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	DockerHubURL     = "https://hub.docker.com"
	DockerHubVersion = "v2"
)

// DockerHubClient represents the data structure of registry servers
type DockerHubClient struct{}

// DockerImage represents an image summery information
type DockerImage struct {
	RepoName         string      `json:"repo_name"`
	ShortDescription string      `json:"short_description"`
	IsOfficial       bool        `json:"is_official"`
	IsAutomated      bool        `json:"is_automated"`
	StarCount        int         `json:"star_count"`
	PullCount        int         `json:"pull_count"`
	RepoOwner        interface{} `json:"repo_owner"`
}

// DockerImageList represents the search results from docker registry server
type DockerImageList struct {
	Previous string        `json:"previous"`
	Next     string        `json:"next"`
	Count    int           `json:"count"`
	Results  []DockerImage `json:"results"`
}

// DockerTag represents docker tag information returned by hub.docker.com
type DockerTag struct {
	Name        string      `json:"name"`
	FullSize    int         `json:"full_size"`
	ID          int64       `json:"id"`
	Repository  int64       `json:"repository"`
	Creator     int64       `json:"creator"`
	LastUpdater int64       `json:"last_updater"`
	LastUpdated string      `json:"last_updated"`
	ImageID     interface{} `json:"image_id"`
	V2          bool        `json:"v2"`
	Platforms   []int       `json:"platforms"`
}

// DockerTagList represents the search results from docker registry server
type DockerTagList struct {
	Previous string      `json:"previous"`
	Next     string      `json:"next"`
	Count    int         `json:"count"`
	Results  []DockerTag `json:"results"`
}

// SendGetRequest sends a request to certain url (basic auth)
func SendGetRequest(url string) (bytes []byte, err error) {
	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, Timeout: 20 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SearchReposByUser returns a list of images in registry
// all image name must begin with repoName+"/"
func (c *DockerHubClient) SearchReposByUser(repoName string) ([]DockerImage, error) {
	if arr := strings.Split(repoName, "/"); len(arr) > 1 {
		return nil, errors.New("only allow repo user passed in")
	}
	if repoName == "" {
		repoName = "library"
	}
	var images []DockerImage

	pageSize := 20
	page := 1
	for {
		var imageList DockerImageList
		url := fmt.Sprintf("%s/%s/search/repositories/?page=%d&query=%s&page_size=%d", DockerHubURL, DockerHubVersion, page, repoName, pageSize)
		bytes, err := SendGetRequest(url)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &imageList)
		if err != nil {
			return nil, err
		}
		glog.V(6).Infof("repo list got, url:%s, count:%d, response:%#v\n", url, len(imageList.Results), imageList.Results)
		for _, image := range imageList.Results {
			if repoName == "library" {
				if image.IsOfficial {
					images = append(images, image)
				}
			} else {
				// only image name starts with repoName+"/", (such as oscarzhao/), it is a qualified result
				if strings.HasPrefix(image.RepoName, repoName+"/") {
					images = append(images, image)
				}
			}
		}

		// next page does not exist, all images got
		if imageList.Next == "" {
			break
		}
		page++
	}
	return images, nil
}

// QueryImageTags returns all tags of certain repo
func (c *DockerHubClient) QueryImageTags(repoName string) ([]DockerTag, error) {
	repoName = strings.Trim(repoName, "/")
	if arr := strings.Split(repoName, "/"); len(arr) == 1 {
		repoName = "library/" + repoName
	}

	var tags []DockerTag
	var tagList DockerTagList
	pageSize := 20
	page := 1
	for {
		url := fmt.Sprintf("%s/%s/repositories/%s/tags/?page=%d&page_size=%d", DockerHubURL, DockerHubVersion, repoName, page, pageSize)
		bytes, err := SendGetRequest(url)
		if err != nil {
			glog.Errorf("fails to fetch tags, url:%s, error:%s\n", url, err)
			return nil, err
		}
		err = json.Unmarshal(bytes, &tagList)
		if err != nil {
			glog.Errorf("invalid response, url:%s, error:%s\n", url, err)
			return nil, err
		}

		tags = append(tags, tagList.Results...)
		if tagList.Next == "" {
			break
		}
		page++
	}
	return tags, nil
}
