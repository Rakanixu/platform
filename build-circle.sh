#!/bin/bash

set -e
set -x

REGISTRY=kazoup


# Used to rebuild all the things

find * -type d -maxdepth 1 -print | while read dir; do
	if [ ! -f $dir/Dockerfile ]; then
		continue
	fi

	pushd $dir >/dev/null


	IMAGE=${dir%/*}-${dir#*/}
	
	# dep
	go get -d  -v -t ./...

	# build static binary
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	

	# build docker image
	docker build -t eu.gcr.io/${PROJECT_NAME}/$IMAGE:$CIRCLE_SHA1 .
	docker tag eu.gcr.io/${PROJECT_NAME}/$IMAGE:$CIRCLE_SHA1 eu.gcr.io/${PROJECT_NAME}/$IMAGE:latest

	
	# remove binary
	rm ${dir#*/}

	popd >/dev/null
done