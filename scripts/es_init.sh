#!/usr/bin/env bash

# Elasticsearch MUST be restarted after this script is executed
sh es_templates/file_template.sh
sh es_indexes/datasource_index.sh
sh es_indexes/flag_index.sh
sh es_plugins/delete_by_query.sh
sh es_plugins/head.sh