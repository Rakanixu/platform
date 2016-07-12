# web-proxy

web-proxy is a microservice acting as proxy for micro web services. Default port to 8000.

FQDN: go.micro.webproxy


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
$ go run main.go
```

