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
        "filename_index_no_ngrams": {
          "type": "custom",
          "filter": [
            "lowercase",
            "asciifolding"
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
    "_default_": {
      "dynamic_templates": [
        {
          "all_text": {
            "match_mapping_type": "string",
            "match": "*",
            "mapping": {
              "type": "string",
              "index": "not_analyzed"
            }
          }
        }
      ]
    },
    "file": {
      "properties": {
        "id": {
          "type": "keyword",
          "index": true
        },
        "user_id": {
          "type": "keyword",
          "index": true
        },
        "name":{
          "type": "text",
          "analyzer": "filename_index_no_ngrams",
          "index": true,
          "fields": {
            "raw": {
              "type": "keyword",
              "index": true
            }
          }
        },
        "category": {
          "type": "keyword",
          "index": true
        },
        "file_type": {
          "type": "keyword",
          "index": true
        },
        "url": {
          "type": "text",
          "analyzer": "path_analyzer",
          "index": true
        },
        "is_dir": {
          "type": "boolean",
          "index": true
        },
        "modified": {
          "type": "date",
          "format": "date_optional_time"
        },
        "last_seen": {
          "type": "long",
          "index": true
        },
        "file_size": {
          "type": "long"
        },
        "access": {
          "type": "keyword",
          "index": true
        },
        "content": {
          "type": "text",
          "analyzer": "content_analyzer",
          "index": true
        },
        "content_category": {
          "type": "keyword",
          "index": true
        },
        "tags": {
          "type": "keyword",
          "index": true
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
