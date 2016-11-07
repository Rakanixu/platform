# Media

Media web is responsible for serving file thumbnails / images previews.

## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service manually

```
$ go run srv/main.go
```

### Run docker containers
Compile Go binaries and build docker image. 
```
make 
```

Run docker container:
```
docker-compose -f docker-compose-build.yml up
```

### Usage
```
/media/preview?user_id=[[md5UserId]]&file_id=[[item.id]]&width=128&height=128&mode=fit&quality=50&token=[[getToken()]]
```
