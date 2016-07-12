# Publish

Publish service with fqdn go.micro.srv.publish
Publish API  with fqdn go.micro.api.publish


## Getting Started

### Prerequisites


### Run Service manually

```
$ go run srv/main.go
```

### Run with NATS

You will need to run gnatsd separately 

```
go run srv/main.go --broker=nats --broker_address=127.0.0.1:4222
```



