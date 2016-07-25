## Usage


### Data

```
go-bindata -pkg data -o data/bindata.go data/

```

### Status
Get Kazoup appliance current status
```
micro query go.micro.srv.config Config.Status '{}'

{
  "APPLIANCE_IS_CONFIGURED": true,
  "APPLIANCE_IS_REGISTERED": true,
  "GIT_COMMIT_STRING": "asdfasdfasdfasdfasdfsdafhash",
  "SMB_USER_EXISTS": true
}
```


### SetElasticSettings
Set ElasticSearch settings for files index. See es_mapping_files.json
```
micro query go.micro.srv.config Config.SetElasticSettings '{}'
{}
```

### SetElasticMapping
Set ElasticSearch mapping for files index and file document type. See es_mapping_files.json
```
micro query go.micro.srv.config Config.SetElasticMapping '{}'
{}
```

### SetFlags
Set kazoup appliances flags into ElasticSearch.. See es_flags.json
```
micro query go.micro.srv.config Config.SetFlags '{}'
{}
```
