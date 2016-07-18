#!/bin/bash

set -e
set -x

REGISTRY=kazoup
go get github.com/mitchellh/gox
# Build Micro
go get github.com/micro/micro
WORKING_DIR=$PWD
# add ES download to build on circle
# https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip

#cd ../../../micro/micro && gox -verbose -os="darwin linux" -arch="386 amd64" -output $WORKING_DIR/bin/micro/{{.OS}}/{{.Arch}}/micro
#cd $WORKING_DIR

#cd ../ui/frontend && npm install && npm install gulp && bower install && node_modules/gulp/bin/gulp.js && cd ../..
# Remove binaries
rm -rf bin
cd ../
find * -type d -maxdepth 1 -print | while read dir; do
	if [ ! -f $dir/Dockerfile ]; then
		continue
	fi

	if [ ${dir%/*} = "ui" ]; then
		continue
	fi
	pushd $dir >/dev/null


	IMAGE=${dir%/*}-${dir#*/}
	
	# dep
	go get -d  -v -t ./...
	
	# test
	go test -v ./...

	# crosscompile
	gox -verbose -os="darwin linux" -arch="386 amd64" -output ../../desktop/bin/${IMAGE}_{{.OS}}_{{.Arch}}
	
	# build ui  		


	popd >/dev/null
done
