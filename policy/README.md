# Policy 

Policy service with fqdn go.micro.srv.policy
Policy API with fqdn go.micro.api.policy


## Getting Started

### Prerequisites


### Run Service manually

```
$ go run srv/main.go
```

### Run API manually

```
$ go run api/main.go
```


### Run docker containers
Compile Go binaries and build docker image. 
```
make 
```

Run docker container:
```
docker-compose -f docker-compose-build.yml up
```


## Usage
[API](https://github.com/kazoup/policy/tree/master/api)

[Microservice](https://github.com/kazoup/policy/tree/master/srv)


