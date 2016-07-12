# Flag API [![GoDoc](https://godoc.org/github.com/Rakanixu/flag/api?status.svg)](https://godoc.org/github.com/Rakanixu/flag/api)

Flag API with fqdn go.micro.api.flag

HTTP interface for flag microservice.


## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run  manually

```
$ go run main.go
```

## Usage

### Create flag
 
```
http[domain:micro API port]/flag/create
{
    "key": "user-advanced-options",
    "description": "Display advanced user options",
    "value": true
}

{}

```


### Read flag
 
```
http[domain:micro API port]/flag/read
{
    "key": "user-advanced-options"
}

{
    "key": "user-advanced-options",
    "description": "Display advanced user options",
    "value": true
}
```


### Flip flag
 
```
http[domain:micro API port]/flag/flip
{
    "key": "user-advanced-options"
}

{}
```


### Delete flag
 
```
http[domain:micro API port]/flag/delete
{
    "key": "user-advanced-options"
}

{}
```


### List flags
 
```
http[domain:micro API port]/flag/list
{}

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


