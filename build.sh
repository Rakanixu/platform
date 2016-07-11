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
	
	# test
	go test -v ./...

	# build static binary
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .

	# build docker image
	docker build -t $REGISTRY/$IMAGE .

	# push docker image
	#docker push $REGISTRY/$IMAGE
	
	# remove binary
	rm ${dir#*/}

	popd >/dev/null
done
