#!/usr/bin/env bash

curl -XPUT localhost:9200/_template/template_file -d '
{
  "template" : "index*",
  "mappings" : {
    "_default_": {
      "dynamic_templates": [
        {
          "all_text": {
            "match_mapping_type": "string",
            "match": "*",
            "mapping": {
              "type": "keyword",
              "index": true
            }
          }
        }
      ]
    },
    "file": {
      "properties": {
        "name":{
          "type": "text",
          "analyzer": "snowball",
          "index": true
        },
        "modified": {
          "type": "date",
          "format": "date_optional_time"
        },
        "content": {
          "type": "text",
          "analyzer": "snowball",
          "index": true,
          "index_options" : "offsets"
        }
      }
    }
  }
}
'
