# Kazoup Analytics Platform  [![CircleCI](https://circleci.com/gh/kazoup/platform/tree/master.svg?style=svg&circle-token=fc062cf6f23c5dc606a8af94b020065a2d073113)](https://circleci.com/gh/kazoup/platform/tree/master)

Kazoup Analytics Platform is build on go-micro library which simplifies developing RPC micro services..

## Services

Docker images for this services are stored and publicly available at [DockerHub](https://hub.docker.com/u/kazoup/)

Repo | Description 
-----|------------  
[meta](https://github.com/kazoup/meta) | Meta data handling service 
[smtp](https://github.com/kazoup/platform/tree/master/smtp) | Email sending service  
[flag](https://github.com/kazoup/platform/tree/master/flag) | Flag micro service 
[policy](https://github.com/kazoup/platfrom/tree/master/policy) | Policy service 
[elastic](https://github.com/kazoup/platform/tree/master/elastic) | Microservice for supporting agnostic CRUD, Search and QueryDSL operations over Elastic search 
[proxy](https://github.com/kazoup/platform/tree/master/proxy) | Entry point for Kazoup platform frontend 
[auth](https://github.com/kazoup/platform/tree/master/auth) | Microservices for Kazoup auth 
[ldap](https://github.com/kazoup/ldap) | LDAP/AD authetication service 
[crawler](https://github.com/kazoup/crawler) | Crawler service 
[publish](https://github.com/kazoup/publish) | Publish service 
[indexer](https://github.com/kazoup/indexer) | Index files from files topic 
[kazoup-web](https://github.com/kazoup/kazoup-web) | Web frontend  as micro services 
## Run

Running everything with Docker is the fastes way to start

```
docker-compose -f docker-compose-all.yml up -d

```




