# Search

Search service, manage all searchs [Local, Network, Google Drive, etc] into one service.

Perform distributed search queries and return back results from different endpoints.

## Getting Started

### Prerequisites
Get Micro
[Micro](https://github.com/micro)
```
go get github.com/micro
```

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


## Usage

[Microservice](https://github.com/Rakanixu/elasticsearch/tree/master/srv)
