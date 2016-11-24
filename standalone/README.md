# A tool to synchronize docker images from gcr.io to index.tenxcloud.com

## Usage 

### Build & Run
```
$ go build -v gcr.go 
$ nohup ./gcr --images=xxx &  
$ # attention: input file path is assigned by flag --images, the format of input file will be talked about in the next section.
```

**Input.** A file with gcr.io/google_containers images, ONE image ONE line, for example:

```
gcr.io/google_containers/pause:2.0
gcr.io/google_containers/pause:3.0  
gcr.io/google_containers/nginx-ingress-controller:0.5
...
```

**Output.** Images in the following list shall be pushed to index.tenxcloud.com/google_containers

```
index.tenxcloud.com/google_containers/pause:2.0
index.tenxcloud.com/google_containers/pause:3.0  
index.tenxcloud.com/google_containers/nginx-ingress-controller:0.5
...
```