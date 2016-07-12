## Usage

### Status
Get kazoup appliance current status
```
[POST] http[domain:micro API port]/config/status
{}

{
  "APPLIANCE_IS_CONFIGURED": true,
  "APPLIANCE_IS_REGISTERED": true,
  "GIT_COMMIT_STRING": "asdfasdfasdfasdfasdfsdafhash",
  "SMB_USER_EXISTS": true
}
```

### SetElasticSettings
Set ElasticSearch settings for files index.
```
[POST] http[domain:micro API port]/config/setElasticSettings
{}

{}
```

### SetElasticMapping
Set ElasticSearch mapping for files index and file document type.
```
[POST] http[domain:micro API port]/config/setElasticMapping
{}

{}
```

### SetFlags
Set kazoup appliances flags into ElasticSearch.
```
[POST] http[domain:micro API port]/config/setFlags
{}

{}
```
