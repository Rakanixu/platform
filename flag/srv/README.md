# Flag service [![GoDoc](https://godoc.org/github.com/Rakanixu/flag/srv?status.svg)](https://godoc.org/github.com/Rakanixu/flag/srv)

Flag microservice with fqdn go.micro.srv.flag

Manage system flags, implementing Create, Read, Flip, Delete and List actions.


## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run manually

```
$ go run main.go
```


## Usage

### Create flag
 
```
micro query go.micro.srv.flag Flag.Create '{"key": "my-unique-flag-key", "description": "You know, for UI feed", "value": true}'
{}
```


### Read flag
 
```
micro query go.micro.srv.flag Flag.Read '{"key": "my-unique-flag-key"}'
{
	"key": "my-unique-flag-key",
	"description": "You know, for UI feed",
	"value": true
}
```


### Flip flag
 
```
micro query go.micro.srv.flag Flag.Flip '{"key": "my-unique-flag-key"}'
{}
```


### Delete flag
 
```
micro query go.micro.srv.flag Flag.Delete '{"key": "my-unique-flag-key"}'
{}
```


### List flags
 
```
micro query go.micro.srv.flag Flag.List '{}'
{
	"result": [
		{
			"key": "my-unique-flag-key-2",
			"description": "You know, for UI feed"
		},
		{
			"key": "my-unique-flag-key-3",
			"description": "You know, for UI feed"
		},
		{
			"key": "my-unique-flag-key-4",
			"description": "You know, for UI feed"
		},
		{
			"key": "my-unique-flag-key-5",
			"description": "You know, for UI feed"
		},
		{
			"key": "my-unique-flag-key-1",
			"description": "You know, for UI feed",
			"value": true
		}
	]
}
```


