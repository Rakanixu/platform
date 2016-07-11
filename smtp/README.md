# SMTP [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go Report Card](https://goreportcard.com/badge/Rakanixu/smtp)](https://goreportcard.com/report/github.com/Rakanixu/smtp)

SMTP microservice and API for email delivery.

## Getting Started

### Prerequisites
Get Micro
[Micro](https://github.com/micro)
```
go get github.com/micro
```

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service manually

```
go run srv/main.go --email_host=SERVER_ADDRESS --email_host_port=SMPT_SERVER_PORT --email_host_user=USERNAME --email_host_password=PASSWORD --default_from_email=noreply@company.com
```

### Run API manually

```
$ go run api/main.go
```


### Run docker containers

Flags have to be passed around, you will want to edit YML file with your SMTP server details.

Compile Go binaries and build docker image. 
```
make 
```

Run docker container:
```
docker-compose -f docker-compose-build.yml up
```


## Usage
[API](https://github.com/Rakanixu/smtp/tree/master/api)

[Microservice](https://github.com/Rakanixu/smtp/tree/master/srv)
