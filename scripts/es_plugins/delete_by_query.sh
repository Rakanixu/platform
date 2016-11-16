#!/usr/bin/env bash

# After installing plugin elasticsearch must be restarted
docker exec -it platform_elasticsearch_1 bin/plugin install delete-by-query
