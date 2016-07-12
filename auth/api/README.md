# Auth API [![Circle CI](https://circleci.com/gh/kazoup/auth-api.svg?style=shield)](https://circleci.com/gh/kazoup/auth-api) [![Go Report Card](https://goreportcard.com/badge/github.com/kazoup/auth-api)](https://goreportcard.com/report/github.com/kazoup/auth-api) [![GoDoc](https://godoc.org/github.com/kazoup/auth-api?status.svg)](https://godoc.org/github.com/kazoup/auth-api)

Auth API is to be used behind Micro API gateway.

API    | ENDPOINT
-------|---------
Auth | /v2/auth/{read}

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```
3. Quick run

	```
	go run main.go
	```
4. Download and start the service

	```shell
	go get github.com/kazoup/auth-api
	auth-api
	```

	OR as a docker container

	```shell
	docker run kazoup/auth-api --registry_address=YOUR_REGISTRY_ADDRESS
	```
