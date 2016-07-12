# INGEST SERVICE [![Circle CI](https://circleci.com/gh/kazoup/kazoup.svg?style=svg&circle-token=47d0513532eacf6f9ce34e18bc1ef22e8325bebc)](https://circleci.com/gh/kazoup/ingest-srv) [![Docker Repository on Quay](https://quay.io/repository/kazoup/ingest-srv/status?token=cdbdfdac-f834-411e-a3bf-cf8159651bb9 "Docker Repository on Quay")](https://quay.io/repository/kazoup/ingest-srv)

ingest-srv is microservice which subscribes to receive file metadata.

future:
Saves metadata into elastic search


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
