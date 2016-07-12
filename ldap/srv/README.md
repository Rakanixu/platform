# LDAP SERVICE

LDAP service is microservice based on go-micro for LDAP/AD binding and authentication.


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
	go get github.com/kazoup/ldap-srv
	ldap-srv
	```

	OR as a docker container

	```shell
	docker run kazoup/ldap-srv --registry_address=YOUR_REGISTRY_ADDRESS
	```
