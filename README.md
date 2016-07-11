# image-sync
synchronize docker images

## build
```
godep go build
```

## run
nohup ./image-sync --src-registry=gcr.io -src-registry-version=v1 --dst-registry=index.tenxcloud.com --dst-registry-version=v2 --dst-repo-password=xxx --repo-owner=google_containers --v=5 &