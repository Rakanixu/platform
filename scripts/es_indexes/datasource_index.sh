#!/usr/bin/env bash

curl -XPOST http://localhost:9200/datasources -d '
{
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
    "datasource": {
      "properties": {
        "user_id": {
          "type": "string",
          "index": "not_analyzed"
        }
      }
    }
  }
}
'
