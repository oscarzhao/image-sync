# image-sync
synchronize docker images

## build
```
godep go build
```

## run
```
# synchronize gcr.io/google_containers to index.tenxcloud.com/google_containers (gcr.io registry version is v1)
docker login index.tenxcloud.com // fill in data
nohup ./image-sync --src-registry=gcr.io -src-registry-version=v1 --dst-registry=index.tenxcloud.com --dst-registry-version=v2 --dst-repo-password=xxx --repo-owner=google_containers --v=5 &

# synchronize index.docker.io/library (official images) to index.tenxcloud.com/docker_library (both registry is version v2)
docker login index.tenxcloud.com (use docker_library as username)
nohup ./image-sync --v=5 &
```
