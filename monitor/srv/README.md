# Monitoring Server

Monitoring server is a microservice which monitors all the applications within the environment. It's a passive system 
used for observation and actions are taken separately based on monitoring events.

The monitoring service subscribes to various events, aggregates and summarises.

```
Who monitors the monitoring service? The monitoring service
```

## How it works

<p align="center">
  <img src="https://github.com/micro/go-os/blob/master/doc/monitor.png" />
</p>

Each service uses [go-os/monitor](https://godoc.org/github.com/micro/go-os/monitor) to create a client 
side Monitor which aggregates Stats, Status and Healthchecks. A user can register healthchecks which are executed 
on a scheduled interval and then published to the monitoring service.

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Download and start the service

	```shell
	go get github.com/micro/monitor-srv
	monitor-srv
	```

	OR as a docker container

	```shell
	docker run microhq/monitor-srv --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Monitoring server implements the following RPC Methods

### Monitor.HealthChecks
```shell
micro query go.micro.srv.monitor Monitor.HealthChecks '{"id": "go.micro.healthcheck.ping"}'
{
	"healthChecks": [
		{
			"description": "This is a ping healthcheck that succeeds",
			"id": "go.micro.healthcheck.ping",
			"results": {
				"foo": "bar",
				"label": "cruft",
				"metric": "1",
				"msg": "a successful check",
				"stats": "123.0"
			},
			"service": {
				"id": "go-server-672d50a7-ab77-11e5-ad84-68a86d0d36b6",
				"name": "go-server",
				"version": "1.0.0"
			},
			"status": 1,
			"timestamp": 1.451096504e+09
		}
	]
}

micro query go.micro.srv.monitor Monitor.HealthChecks '{"id": "go.micro.healthcheck.pong"}'
{
	"healthChecks": [
		{
			"description": "This is a pong healthcheck that fails",
			"error": "Unknown exception",
			"id": "go.micro.healthcheck.pong",
			"results": {
				"foo": "ugh",
				"label": "",
				"metric": "-0.0001",
				"msg": "a catastrophic failure occurred",
				"stats": "NaN"
			},
			"service": {
				"id": "go-server-672d50a7-ab77-11e5-ad84-68a86d0d36b6",
				"name": "go-server",
				"version": "1.0.0"
			},
			"status": 2,
			"timestamp": 1.451096508e+09
		}
	]
}
```

### Sending HealthChecks

Healthchecks are sent using [go-os/monitor](https://github.com/micro/go-os/blob/master/monitor/monitor.go). Example usage can be found in [go-os/examples/monitor](https://github.com/micro/go-os/blob/master/examples/monitor/monitor.go).

