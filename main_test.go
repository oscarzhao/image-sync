package main

import (
	"testing"
)

func TestImageStruct(t *testing.T) {
	img := Image{
		registry: "gcr.io",
		repo:     "google_containers/ubuntu",
		tag:      "14.04",
	}
	fullName := "gcr.io/google_containers/ubuntu:14.04"

	if img.String() != fullName {
		t.Fatalf("should be %s, is %s\n", fullName, img)
	}
}
