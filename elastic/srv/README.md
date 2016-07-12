# Elasticsearch service [![GoDoc](https://godoc.org/github.com/Rakanixu/elasticsearch/srv?status.svg)](https://godoc.org/github.com/Rakanixu/elasticsearch/srv)

This is the elasticsearch service performing CRUD, Search and QueryDSL operations for Elasticsearch DB.

## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

### Run Service

```
$ go run main.go
```

## Usage 

### Create
```
micro query go.micro.srv.elastic Elastic.Create '{"index":"flags", "type": "flag", "id": "flag-id", "data":  "{\"att1\": \"value1\", \"bool\": false, \"innerobj\":{\"attr1\": 46,\"bool\": true}}"}'
{}
```

### Read
```
micro query go.micro.srv.elastic Elastic.Read '{"index":"flags", "type": "flag", "id": "flag-id"}'
{
	"result": "{\"att1\": \"value1\", \"bool\": false, \"innerobj\":{\"attr1\": 46,\"bool\": true}}"
}

```

### Update
```
micro query go.micro.srv.elastic Elastic.Update '{"index":"flags", "type": "flag", "id": "flag-id", "data":  "{\"att1\": \"new value1\", \"bool\": false, \"innerobj\":{\"attr1\": 46,\"bool\": true}}"}'
{}
```


### Delete
```
micro query go.micro.srv.elastic Elastic.Delete '{"index":"flags", "type": "flag", "id": "flag-id"}'
{}
```

### Search
```
micro query go.micro.srv.elastic Elastic.Search '{"index":"flags", "type": "flag", "query": "yy", "limit": 20, "offset": 0}'
{
	"result": "{\"took\":1,\"timed_out\":false,\"_shards\":{\"total\":5,\"successful\":5,\"failed\":0},\"hits\":{\"total\":4,\"max_score\":1.0,\"hits\":[{\"_index\":\"flags\",\"_type\":\"flag\",\"_id\":\"flag-27\",\"_score\":1.0,\"_source\":{\"fieldY\": \"yy\", \"bb\": false, \"obj\":{\"obj2\": 46}}},{\"_index\":\"flags\",\"_type\":\"flag\",\"_id\":\"flag-25\",\"_score\":1.0,\"_source\":{\"fieldY\": \"aa\", \"bb\": true, \"obj\":{\"obj1\": 44}}},{\"_index\":\"flags\",\"_type\":\"flag\",\"_id\":\"flag-31\",\"_score\":1.0,\"_source\":{\"fieldY\": \"yyt\", \"bb\": false, \"obj\":{\"obj2\": 46,\"obj55\": 66666666666666}}},{\"_index\":\"flags\",\"_type\":\"flag\",\"_id\":\"flag-30\",\"_score\":1.0,\"_source\":{\"fieldY\": \"yy\", \"bb\": false, \"obj\":{\"obj2\": 46,\"obj1\": 6666666666666666666}}}]}}"
}
```

### Query
```
micro query go.micro.srv.elastic Elastic.Query '{"index":"flags", "type": "flag", "query": "{\"query\":{\"match\" : {\"fieldY\" : \"yy\"}}}"}'
{
	"result": "{\"took\":2,\"timed_out\":false,\"_shards\":{\"total\":5,\"successful\":5,\"failed\":0},\"hits\":{\"total\":2,\"max_score\":0.30685282,\"hits\":[{\"_index\":\"flags\",\"_type\":\"flag\",\"_id\":\"flag-27\",\"_score\":0.30685282,\"_source\":{\"fieldY\": \"yy\", \"bb\": false, \"obj\":{\"obj2\": 46}}},{\"_index\":\"flags\",\"_type\":\"flag\",\"_id\":\"flag-30\",\"_score\":0.30685282,\"_source\":{\"fieldY\": \"yy\", \"bb\": false, \"obj\":{\"obj2\": 46,\"obj1\": 6666666666666666666}}}]}}"
}

```


### CreateIndexWithSettings
Note: scaping character "\"
eg) \"pattern\": \"[^\\\\p{L}\\\\d]+\"
For reference, compare this call with the API call you can find in API Readme. Both works, but API is human friendly.
```
micro query go.micro.srv.elastic Elastic.CreateIndexWithSettings '{"index": "filestest", "settings":"{\"settings\": {\"analysis\": {\"analyzer\": { \"split_on_bar\": {\"type\": \"custom\",\"tokenizer\": \"split_on_bar\"}, \"filename_index\": { \"type\": \"custom\", \"filter\": [ \"lowercase\", \"asciifolding\", \"ngram_3_20\" ], \"tokenizer\": \"filename\" } }, \"filter\": { \"ngram_3_20\": { \"type\": \"nGram\", \"min_gram\": \"3\", \"max_gram\": \"20\" } }, \"tokenizer\": { \"split_on_bar\": { \"pattern\": \"[|]\", \"type\": \"pattern\" }, \"filename\": { \"pattern\": \"[^\\\\p{L}\\\\d]+\", \"type\": \"pattern\" } } } } }"}'
{}
```

### PutMappingFromJSON
```
micro query go.micro.srv.elastic Elastic.PutMappingFromJSON '{"index": "files","type": "file","mapping":"{ESCAPED_STRINGiFY_JSON}"}'
{}
```
