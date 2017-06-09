# AGENT SERVICE

agent-srv is microservice that is used for creating Kazoup files depending on the
data sent from the client installed on desktop pc / server.

## Getting started

### Prerequisites

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
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
docker-compose -f docker-compose-build.yml up agent
```

### Usage

TODO: Create example request
