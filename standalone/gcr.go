package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

var (
	imageFile string
	src       = "gcr.io"
	dst       = "index.tenxcloud.com"
)

func init() {
	flag.StringVar(&imageFile, "images", "gcr.io", "a file storing a list of docker images starting with gcr.io/google_containers")
}

func main() {
	rawBytes, err := ioutil.ReadFile(imageFile)
	if err != nil {
		log.Fatalf("read file %s fails, error:%v\n", imageFile, err)
	}

	arr := strings.Split(string(rawBytes), "\n")
	var images []string
	for _, img := range arr {
		img = strings.Trim(img, "\n\t\r ")
		if len(img) > 0 && strings.HasPrefix(img, "gcr.io/google_containers") {
			images = append(images, img)
		}
	}
	if len(images) == 0 {
		log.Fatalf("invalid config file, no gcr.io/google_containers images\n")
	}

	for _, img := range images {
		stdout, stderr, err := cmdExec([]string{"docker", "pull", img})
		if err != nil {
			log.Printf("pull image %s fails, error:%s, stdout:%s, stderr:%s\n", img, err, stdout, stderr)
			continue
		}

		dstImg := strings.Replace(img, src, dst, -1)
		stdout, stderr, err = cmdExec([]string{"docker", "tag", img, dstImg})
		if err != nil {
			log.Printf("tag image %s to %s fails, error:%s, stdout:%s, stderr:%s\n", img, dstImg, err, stdout, stderr)
			continue
		}

		stdout, stderr, err = cmdExec([]string{"docker", "push", dstImg})
		if err != nil {
			log.Printf("push image %s fails, error:%s, stdout:%s, stderr:%s\n", dstImg, err, stdout, stderr)
			continue
		}

		stdout, stderr, err = cmdExec([]string{"docker", "rmi", dstImg})
		if err != nil {
			log.Printf("rm image %s fails, error:%s, stdout:%s, stderr:%s\n", dstImg, err, stdout, stderr)
			continue
		}

		stdout, stderr, err = cmdExec([]string{"docker", "rmi", img})
		if err != nil {
			log.Printf("rm image %s fails, error:%s, stdout:%s, stderr:%s\n", img, err, stdout, stderr)
			continue
		}
	}
}

func cmdExec(cmds []string) ([]byte, []byte, error) {
	var stdoutB, stderrB bytes.Buffer
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB
	err := cmd.Run()
	return stdoutB.Bytes(), stderrB.Bytes(), err
}
