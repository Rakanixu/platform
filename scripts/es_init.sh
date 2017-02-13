#!/usr/bin/env bash

# Elasticsearch MUST be restarted after this script is executed
sh es_templates/file_template.sh
sh es_indexes/datasource_index.sh
