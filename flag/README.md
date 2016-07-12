# Flag [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go Report Card](https://goreportcard.com/badge/Rakanixu/flag)](https://goreportcard.com/report/github.com/Rakanixu/flag)

Flag service with fqdn go.micro.srv.flag
Flag API with fqdn go.micro.api.flag

Data is stored in Elasticsearch.


## Getting Started

### Prerequisites
Get Micro
[Micro](https://github.com/micro)
```
go get github.com/micro
```

This microservice needs elatiscsearch service up and running. 
```
go get github.com/Rakanixu/elasticsearch
make
docker-compose -f docker-compose-build.yml up
```
Now we've got Consul, Elasticsearch (DB), micro api, micro web, elasticsearch-srv and elasticsearch-api up and running. 


### Run Service manually

```
$ go run srv/main.go
```

### Run API manually

```
$ go run api/main.go
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
[API](https://github.com/Rakanixu/flag/tree/master/api)

[Microservice](https://github.com/Rakanixu/flag/tree/master/srv)


