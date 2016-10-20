# CRAWLER SERVICE

crawler-srv is microservice to publish file metadata.

## Crawlers
- fakescanner: Generates random data


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
	go get github.com/kazoup/crawler-srv
	ldap-srv
	```

 	OR as a docker container

 	```shell
 	docker run kazoup/crawler-srv --registry_address=YOUR_REGISTRY_ADDRESS
 	```
