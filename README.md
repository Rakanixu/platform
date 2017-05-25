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

## Run Development Environment

1.  Update you hosts file (/etc/hosts on UNIX-like OS)
    Add following line to resolve those names to localhost:
    ```
    127.0.0.1	web.kazoup.io elasticsearch app.kazoup.io nats tika redis kibana
    ```

2.  Run platform and dependencies with docker compose.
    In platform directory there are  two YML files, services.yml and platform.yml
    
    ```
    cd /$GOPATH/src/github.com/kazoup/platform
    
    docker-compose -f services.yml up
    ```
    
    When stoping services, registry container gets an unhealthy state. Remove all data from registry container and restart it (Just registry or all containers all together)
    When removing elasticsearch container, all data will be lost, so be sure to execute step 3 every time elasticsearch is wiped out.
    
    ```
    docker-compose -f services.yml rm registry
    ```
    
3.  Elasticsearch is running now and should be accesible on localhost:9200, but we need to apply Kazoup settings and mappings.
    
    ```
    ./scripts/es_system/sysctl.sh
    
    ./scripts/es_init.sh
    ```
    
4.  Run the Kazoup platform microservices:

    ```
    docker-compose -f platform.yml up
    ```

5.  Serve frontend SPA:

    ```
    cd /$GOPATH/src/github.com/kazoup/platform/ui/polymer
    
    polymer serve
    ```

6. Check platform status on the browser and app
    ```
    https://localhost:8082/
    
    https://localhost:8082/registry
    
    http://localhost:8080
    ```

Kazoup source code is ussually not compiled when running the development environment, (have a look on platform.yml)
That's why it is recommendable to build the service you are working on or whole platform before running it.
```
go build ./...
```

Before pushing

```
go test ./...
```

When updating protocol buffers, regenate them by 
```
cd /$GOPATH/src/github.com/kazoup/platform/SERVICE_NAME/
make protoc
```

