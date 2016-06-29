# Kazoup Analytics Platform 

Kazoup Analytics Platform is build on go-micro library which simplifies developing RPC micro services.

## Services

Docker images for this services are stored and pulicly available at [DockerHub](https://hub.docker.com/u/kazoup/)

Repo | Description |   CI     
-----|------------ | -------- 
[micro](https://github.com/micro/micro) |Micro services toolkit | [![Build Status](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) 
[meta](https://github.com/kazoup/meta) | Meta data handling service | [![CircleCI](https://circleci.com/gh/kazoup/meta.svg?style=svg)](https://circleci.com/gh/kazoup/meta)
[smtp](https://github.com/kazoup/smtp) | SMTP microservice for mail delivery | [![CircleCI](https://circleci.com/gh/kazoup/smtp.svg?style=svg)](https://circleci.com/gh/kazoup/smtp)
[flag](https://github.com/kazoup/flag) | Flag micro service | [![CircleCI](https://circleci.com/gh/kazoup/flag.svg?style=svg)](https://circleci.com/gh/kazoup/flag) 
[policy](https://github.com/kazoup/policy) | Policy service |
[elastic](https://github.com/kazoup/elastic) | Microservice for supporting agnostic CRUD, Search and QueryDSL operations over Elastic search | [![CircleCI](https://circleci.com/gh/kazoup/elastic.svg?style=svg)](https://circleci.com/gh/kazoup/elastic)




## Run

Running everything with Docker is the fastes way to start

```
docker-compose -f docker-compose-prod.yml up -d

```

## Build


```

make 

```



