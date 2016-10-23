#!/usr/bin/env bash

curl -XPUT localhost:9200/_template/template_file -d '
{
  "template" : "index*",
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
        },
        "path_analyzer": {
          "type": "custom",
          "tokenizer": "path_tokenizer"
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
        },
        "path_tokenizer": {
          "type": "path_hierarchy"
        }
      }
    }
  },
  "mappings" : {
    "file": {
      "properties": {
        "id": {
          "type": "string",
          "index": "not_analyzed"
        },
        "user_id": {
          "type": "string",
          "index": "not_analyzed"
        },
        "name":{
          "type": "string",
          "analyzer": "filename_index",
          "fields": {
            "raw": {
              "type": "string",
              "index": "not_analyzed"
            }
          }
        },
        "category": {
          "type": "string",
          "index": "not_analyzed"
        },
        "url": {
          "type": "string",
          "analyzer": "path_analyzer"
        },
        "is_dir": {
          "type": "boolean"
        },
        "modified": {
          "type": "date",
          "format": "date_optional_time"
        },
        "file_size": {
          "type": "long"
        }
      }
    }
  }
}
'
