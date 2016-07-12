# Elasticsearch API [![GoDoc](https://godoc.org/github.com/Rakanixu/elasticsearch/api?status.svg)](https://godoc.org/github.com/Rakanixu/elasticsearch/api)

This is the Elasticsearch API for consuming elascticsearch service through HTTP.

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
http[domain:micro API port]/elastic/create
{
    "index":"flags", 
    "type": "flag", 
    "id": "flag-id", 
    "data":  {
        "att1": "value1", 
        "bool": false, 
        "innerobj": {
            "attr1": 46,
            "bool": true
        }
    }
}

{}
```

### Read
```
http[domain:micro API port]/elastic/read
{
    "index":"flags", 
    "type": "flag", 
    "id": "flag-id"
}

{
  "att1": "value1",
  "bool": false,
  "innerobj": {
    "attr1": 46,
    "bool": true
  }
}
```

### Update
```
http[domain:micro API port]/elastic/update
{
    "index":"flags", 
    "type": "flag", 
    "id": "flag-id",
    "data":  {
        "update": true,
        "att1": "value1", 
        "bool": false, 
        "innerobj": {
            "attr1": 46,
            "bool": true
        }
    }
}

{
    "update": true,
    "att1": "value1", 
    "bool": false, 
    "innerobj": {
        "attr1": 46,
        "bool": true
    }
}
```

### Delete
```
http[domain:micro API port]/elastic/delete
{
    "index":"flags", 
    "type": "flag", 
    "id": "flag-id"
}

{}
```

### Search
```
http[domain:micro API port]/elastic/search
{
    "index":"flags", 
    "type": "flag", 
    "query":"47", 
    "limit": 20, 
    "offset": 0
}

{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 0.19178301,
    "hits": [
      {
        "_index": "flags",
        "_type": "flag",
        "_id": "flag-id3",
        "_score": 0.19178301,
        "_source": {
          "att1": "value2",
          "bool": false,
          "innerobj": {
            "attr1": 47,
            "bool": true
          }
        }
      }
    ]
  }
}
```

### Query
```
http[domain:micro API port]/elastic/query
{
    "index":"flags", 
    "type": "flag", 
    "query": {
        "query": { 
            "match" : {
                "att1" : "value2"
            }
        }
    }
}

{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 2,
    "max_score": 0.30685282,
    "hits": [
      {
        "_index": "flags",
        "_type": "flag",
        "_id": "flag-id2",
        "_score": 0.30685282,
        "_source": {
          "att1": "value2",
          "bool": false,
          "innerobj": {
            "attr1": 48,
            "bool": true
          },
          "update": true
        }
      },
      {
        "_index": "flags",
        "_type": "flag",
        "_id": "flag-id3",
        "_score": 0.30685282,
        "_source": {
          "att1": "value2",
          "bool": false,
          "innerobj": {
            "attr1": 47,
            "bool": true
          }
        }
      }
    ]
  }
}
```

### CreateIndexWithSettings
```
http[domain:micro API port]/elastic/createIndexWithSettings
{
    "index": "files",
    "settings": {
      "settings": {
        "analysis": {
          "analyzer": {
            "split_on_bar": {
              "type": "custom",
              "tokenizer": "split_on_bar"
            },
            "filename_index": {
              "type": "custom",
              "filter": [
                "lowercase",
                "asciifolding",
                "ngram_3_20"
              ],
              "tokenizer": "filename"
            }
          },
          "filter": {
            "ngram_3_20": {
              "type": "nGram",
              "min_gram": "3",
              "max_gram": "20"
            }
          },
          "tokenizer": {
            "split_on_bar": {
              "pattern": "[|]",
              "type": "pattern"
            },
            "filename": {
              "pattern": "[^\\p{L}\\d]+",
              "type": "pattern"
            }
          }
        }
      }
    }
}

{}
```

### PutMappingFromJSON
```
http[domain:micro API port]/elastic/putMappingFromJSON
{  
   "index":"files",
   "type":"file",
   "mapping":{  
      "dynamic_templates":[  
         {  
            "date_fields":{  
               "mapping":{  
                  "doc_values":true,
                  "format":"date_optional_time",
                  "type":"date"
               },
               "match":".*",
               "match_mapping_type":"date",
               "match_pattern":"regex"
            }
         },
         {  
            "default_not_analyzed_with_doc_values":{  
               "mapping":{  
                  "index":"not_analyzed",
                  "doc_values":true,
                  "type":"{dynamic_type}"
               },
               "match":".*",
               "match_mapping_type":"string|boolean|double|long|integer",
               "match_pattern":"regex"
            }
         }
      ],
      "properties":{  
         "content":{  
            "properties":{  
               "checksum":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "checksum_time":{  
                  "type":"double",
                  "doc_values":true
               },
               "content_stored":{  
                  "type":"boolean"
               },
               "content_time":{  
                  "type":"double",
                  "doc_values":true
               },
               "date":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "location":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "money":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "organisation":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "percent":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "person":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "text":{  
                  "type":"string",
                  "analyzer":"standard"
               },
               "tika_metadata":{  
                  "type":"string",
                  "analyzer":"standard"
               },
               "time":{  
                  "properties":{  
                     "count":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "value":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "unique":{  
                  "type":"boolean"
               }
            }
         },
         "exists_on_disk":{  
            "type":"boolean"
         },
         "first_seen":{  
            "type":"date",
            "doc_values":true,
            "format":"date_optional_time"
         },
         "id_b64":{  
            "type":"string",
            "index":"not_analyzed",
            "doc_values":true
         },
         "last_seen":{  
            "type":"date",
            "doc_values":true,
            "format":"date_optional_time"
         },
         "metadata":{  
            "properties":{  
               "accessed":{  
                  "type":"date",
                  "doc_values":true,
                  "format":"date_optional_time"
               },
               "created":{  
                  "type":"date",
                  "doc_values":true,
                  "format":"date_optional_time"
               },
               "dirpath":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "dirpath_split":{  
                  "properties":{  
                     "0":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "1":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "10":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "11":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "12":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "13":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "14":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "15":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "16":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "17":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "18":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "19":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "2":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "20":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "3":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "4":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "5":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "6":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "7":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "8":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "9":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "doc_type":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "extension":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "filename":{  
                  "type":"string",
                  "analyzer":"filename_index",
                  "fields":{  
                     "raw":{  
                        "type":"string",
                        "index":"not_analyzed",
                        "doc_values":true
                     }
                  }
               },
               "filename_b64":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "fullpath":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "gid":{  
                  "type":"long",
                  "doc_values":true
               },
               "mimetype":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "modified":{  
                  "type":"date",
                  "doc_values":true,
                  "format":"date_optional_time"
               },
               "sharepath":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "size":{  
                  "type":"long",
                  "doc_values":true
               },
               "smb_attributes":{  
                  "properties":{  
                     "created":{  
                        "type":"long",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "extended_attributes":{  
                        "type":"long",
                        "doc_values":true
                     },
                     "last_access":{  
                        "type":"long",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "last_attr_change":{  
                        "type":"long",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "last_write":{  
                        "type":"long",
                        "index":"not_analyzed",
                        "doc_values":true
                     },
                     "offline":{  
                        "type":"boolean"
                     }
                  }
               },
               "uid":{  
                  "type":"long",
                  "doc_values":true
               }
            }
         },
         "permissions":{  
            "properties":{  
               "access_groups":{  
                  "type":"string",
                  "analyzer":"split_on_bar"
               },
               "access_users":{  
                  "type":"string",
                  "analyzer":"split_on_bar"
               },
               "acl":{  
                  "type":"string",
                  "analyzer":"split_on_bar"
               },
               "acl_error":{  
                  "type":"string",
                  "index":"not_analyzed",
                  "doc_values":true
               },
               "allow":{  
                  "type":"string",
                  "analyzer":"split_on_bar"
               },
               "deny":{  
                  "type":"string",
                  "analyzer":"split_on_bar"
               }
            }
         },
         "unique":{  
            "type":"boolean"
         }
      }
   }
}

{}
```
