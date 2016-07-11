# Kazoup Analytics Platform 

Kazoup Analytics Platform is build on go-micro library which simplifies developing RPC micro services.

## Services

Docker images for this services are stored and publicly available at [DockerHub](https://hub.docker.com/u/kazoup/)

Repo | Description 
-----|------------  
[meta](https://github.com/kazoup/meta) | Meta data handling service 
[smtp](https://github.com/kazoup/platform/tree/master/smtp) | Email sending service  
[flag](https://github.com/kazoup/flag) | Flag micro service 
[policy](https://github.com/kazoup/policy) | Policy service 
[elastic](https://github.com/kazoup/elastic) | Microservice for supporting agnostic CRUD, Search and QueryDSL operations over Elastic search 
[web-proxy](https://github.com/kazoup/web-proxy) | Entry point for Kazoup platform frontend 
[auth](https://github.com/kazoup/auth) | Microservices for Kazoup auth 
[kazoup-web](https://github.com/kazoup/kazoup-web) | Web frontend  as micro services 
[ldap](https://github.com/kazoup/ldap) | LDAP/AD authetication service 
[crawler](https://github.com/kazoup/crawler) | Crawler service 
[publish](https://github.com/kazoup/publish) | Publish service 
[indexer](https://github.com/kazoup/indexer) | Index files from files topic 
## Run

Running everything with Docker is the fastes way to start

```
docker-compose -f docker-compose-all.yml up -d

```




