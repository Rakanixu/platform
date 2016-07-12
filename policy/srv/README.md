# Policy srv

Policy srv for managing Kazoup data policies.
A policy can never be updated, only created and deleted.

FQDN: go.micro.srv.policy

## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service

```
$ go run main.go
```

### Usage

#### Create
```
micro query go.micro.srv.policy Policy.Create '{"filter": "{\"bool\":{\"must\":[{\"prefix\":{\"metadata.dirpath\":\"//10.17.57.20/print$/color\"}},{\"range\":{\"accessed\":{\"gte\":\"now-26y/d\",\"lte\":\"now-0y/d\"}}}]}}", "filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D", "is_archive_policy": true, "is_deletion_policy": false, "name": "policy name"}'
{}
```

#### Read
```
micro query go.micro.srv.policy Policy.Read '{"name": "policy name"}'
{
	"created": "2016-06-28 14:42:20.899444558 +0100 BST",
	"created_by": "unknown",
	"filter": "{\"bool\":{\"must\":[{\"prefix\":{\"metadata.dirpath\":\"//10.17.57.20/print$/color\"}},{\"range\":{\"accessed\":{\"gte\":\"now-26y/d\",\"lte\":\"now-0y/d\"}}}]}}",
	"filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
	"is_archive_policy": true,
	"name": "policy name"
}
```

#### Delete
```
micro query go.micro.srv.policy Policy.Delete '{"name": "policy name"}'
{}
```


#### List
```
micro query go.micro.srv.policy Policy.List '{}'
{
	"result": [
		{
			"created": "2016-06-28 14:44:55.116953655 +0100 BST",
			"created_by": "unknown",
			"filter": "{\"bool\":{\"must\":[{\"prefix\":{\"metadata.dirpath\":\"//10.17.57.20/print$/color\"}},{\"range\":{\"accessed\":{\"gte\":\"now-26y/d\",\"lte\":\"now-0y/d\"}}}]}}",
			"filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
			"name": "policy name 2"
		},
		{
			"created": "2016-06-28 14:42:20.899444558 +0100 BST",
			"created_by": "unknown",
			"filter": "{\"bool\":{\"must\":[{\"prefix\":{\"metadata.dirpath\":\"//10.17.57.20/print$/color\"}},{\"range\":{\"accessed\":{\"gte\":\"now-26y/d\",\"lte\":\"now-0y/d\"}}}]}}",
			"filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
			"is_archive_policy": true,
			"name": "policy name"
		},
		{
			"created": "2016-06-28 14:44:16.94757734 +0100 BST",
			"created_by": "unknown",
			"filter": "{\"bool\":{\"must\":[{\"prefix\":{\"metadata.dirpath\":\"//10.17.57.20/print$/color\"}},{\"range\":{\"accessed\":{\"gte\":\"now-26y/d\",\"lte\":\"now-0y/d\"}}}]}}",
			"filter_raw": "%7B%22from%22%3A%22now-26y%2Fd%22%2C%22to%22%3A%22now-0y%2Fd%22%2C%22dirpath%22%3A%5B%22%2F%2F10.17.57.20%2Fprint%24%2Fcolor%22%5D%7D",
			"is_archive_policy": true,
			"name": "policy name 1"
		}
	]
}
```
