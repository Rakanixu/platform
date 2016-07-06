# Kazoup Analytics Platform 

Kazoup Analytics Platform is build on go-micro library which simplifies developing RPC micro services.

## Services

Docker images for this services are stored and publicly available at [DockerHub](https://hub.docker.com/u/kazoup/)

Repo | Description |   CI     
-----|------------ | -------- 
[micro](https://github.com/micro/micro) |Micro services toolkit | [![Build Status](https://travis-ci.org/micro/micro.svg?branch=master)](https://travis-ci.org/micro/micro) 
[meta](https://github.com/kazoup/meta) | Meta data handling service | [![CircleCI](https://circleci.com/gh/kazoup/meta.svg?style=svg)](https://circleci.com/gh/kazoup/meta)
[smtp](https://github.com/kazoup/smtp) | SMTP microservice for mail delivery | [![CircleCI](https://circleci.com/gh/kazoup/smtp.svg?style=svg)](https://circleci.com/gh/kazoup/smtp)
[flag](https://github.com/kazoup/flag) | Flag micro service | [![CircleCI](https://circleci.com/gh/kazoup/flag.svg?style=svg)](https://circleci.com/gh/kazoup/flag) 
[policy](https://github.com/kazoup/policy) | Policy service | [![CircleCI](https://circleci.com/gh/kazoup/policy.svg?style=svg&circle-token=1e5f2d34488ed3bad550549f76e6ec45eca6c50d)](https://circleci.com/gh/kazoup/policy)
[elastic](https://github.com/kazoup/elastic) | Microservice for supporting agnostic CRUD, Search and QueryDSL operations over Elastic search | [![CircleCI](https://circleci.com/gh/kazoup/elastic.svg?style=svg)](https://circleci.com/gh/kazoup/elastic)
[web-proxy](https://github.com/kazoup/web-proxy) | Entry point for Kazoup platform frontend | [![CircleCI](https://circleci.com/gh/kazoup/web-proxy.svg?style=svg&circle-token=1644b35cf078b8382f46748e39299d525ce15fc0)](https://circleci.com/gh/kazoup/web-proxy)
[auth](https://github.com/kazoup/auth) | Microservices for Kazoup auth | [![CircleCI](https://circleci.com/gh/kazoup/auth.svg?style=svg&circle-token=fb3082b3ae297e36628bbba40e69eb0a3d8fe247)](https://circleci.com/gh/kazoup/auth)
[kazoup-web](https://github.com/kazoup/kazoup-web) | Web frontend  as micro services | [![CircleCI](https://circleci.com/gh/kazoup/kazoup-web.svg?style=svg&circle-token=1084085b649711ccdac2e6355412dcd9fb259f64)](https://circleci.com/gh/kazoup/kazoup-web)
[ldap](https://github.com/kazoup/ldap) | LDAP/AD authetication service | [![CircleCI](https://circleci.com/gh/kazoup/ldap.svg?style=svg&circle-token=c5e2408d51b764c10b2736213c754339996feee1)](https://circleci.com/gh/kazoup/ldap)



## Run

Running everything with Docker is the fastes way to start

```
docker-compose -f docker-compose-all.yml up -d

```

## Build


```

make 

```



