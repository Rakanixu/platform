# Kazoup Platform  [![CircleCI](https://circleci.com/gh/kazoup/platform/tree/master.svg?style=svg&circle-token=fc062cf6f23c5dc606a8af94b020065a2d073113)](https://circleci.com/gh/kazoup/platform/tree/master)

Kazoup  Platform is build on go-micro library which simplifies developing RPC micro services..




## Services

Docker images for this services are stored and publicly available at [DockerHub](https://hub.docker.com/u/kazoup/)



## Test

This will trun all tests

```
./test.sh

```



## Build

This will build all platform 

```
./build.sh

```

## Deploy

Deploy images to DockerHub

```

./deploy.sh

```

## Run

Running everything with Docker is the fastes way to start

Set your /etc/hosts file to 

```
127.0.0.1 web.kazoup.io elasticsearch app.kazoup.io

```

```
docker-compose -f docker-compose-all.yml up -d

```




