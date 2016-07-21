#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
# Build Micro
go get github.com/micro/micro
WORKING_DIR=$PWD
# add ES download to build on circle
# https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip

#cd ../../../micro/micro && gox -verbose -os="darwin linux" -arch="386 amd64" -output $WORKING_DIR/bin/micro/{{.OS}}/{{.Arch}}/micro
#cd $WORKING_DIR

#cd ../ui/frontend && npm install && npm install gulp && bower install && node_modules/gulp/bin/gulp.js && cd ../..

# Build binary distributable
rm -rf desktop/bin



	
	# dep
	
	# test
	go test -v ./...

	# crosscompile
	gox -verbose -os="darwin linux" -arch="386 amd64" -output bin/{{.OS}}/{{.Arch}}/kazoup-{{.OS}}-{{.Arch}}
	
	# build ui  		

