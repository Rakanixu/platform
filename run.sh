#!/bin/bash

set -e
set -x

REGISTRY=kazoup

micro --registry=mdns api &
micro --registry=mdns web &

# Used to rebuild all the things

find * -type d -maxdepth 1 -print | while read dir; do
	if [ ! -f $dir/Dockerfile ]; then
		continue
	fi

	pushd $dir >/dev/null


	IMAGE=${dir%/*}-${dir#*/}

	if [[ ${dir%/*} = "elastic" && ${dir#*/} = "srv" ]]; then
  		go run main.go --registry=mdns --elastic_hosts=127.0.0.1:9200 &
  		break
  fi
	go run main.go --registry=mdns &

	popd >/dev/null
done
