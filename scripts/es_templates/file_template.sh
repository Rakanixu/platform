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
        "content_analyzer": {
          "type": "custom",
          "char_filter" : [ "html_strip" ],
          "filter": [
            "lowercase",
            "asciifolding"
          ],
          "tokenizer": "content"
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
        "content": {
          "type": "whitespace"
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
        "file_type": {
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
        "last_seen": {
          "type": "long",
          "index": "not_analyzed"
        },
        "file_size": {
          "type": "long"
        },
        "access": {
          "type": "string",
          "index": "not_analyzed"
        },
        "content": {
          "type": "string",
          "analyzer": "content_analyzer"
        },
        "content_category": {
          "type": "string",
          "index": "not_analyzed"
        },
        "tags": {
          "type": "string",
          "index": "not_analyzed"
        },
        "opts_kazoup_file": {
          "type": "nested",
          "properties": {
            "tags_timestamp": {
              "type": "date",
              "format": "date_optional_time"
            },
            "content_timestamp": {
              "type": "date",
              "format": "date_optional_time"
            },
            "audio_timestamp": {
              "type": "date",
              "format": "date_optional_time"
            },
            "text_analyzed_timestamp": {
              "type": "date",
              "format": "date_optional_time"
            }
          }
        }
      }
    }
  }
}
'
