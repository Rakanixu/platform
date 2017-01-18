# Monitor Web

The monitor web is a dashboard for the monitor service. 

## Getting started

1. Install Consul

	Consul is the default registry/monitor for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -dev -advertise=127.0.0.1
	```

3. Download and start the service

	```shell
	go get github.com/micro/monitor-web
	monitor-web
	```

	OR as a docker container

	```shell
	docker run microhq/monitor-web --registry_address=YOUR_REGISTRY_ADDRESS
	```

![Monitor Web 1](image1.png)
-
![Monitor Web 2](image2.png)
-
![Monitor Web 3](image3.png)

