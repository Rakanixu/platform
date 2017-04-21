# Channel

Channel service responsible for retrieving channels related with a file (metadata)

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

