#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
# Build Micro
go get github.com/micro/micro
WORKING_DIR=$PWD
DESKTOP_DIR=$GOPATH/github.com/kazoup/platform/desktop
# add ES download to build on circle
# https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip

#cd ../../../micro/micro && gox -verbose -os="darwin linux" -arch="386 amd64" -output $WORKING_DIR/bin/micro/{{.OS}}/{{.Arch}}/micro
#cd $WORKING_DIR

#cd ../ui/frontend && npm install && npm install gulp && bower install && node_modules/gulp/bin/gulp.js && cd ../..

# Build binary distributable
rm -rf  $GOPATH/src/github.com/kazoup/platform/desktop/bin
# crosscompile Kazoup CLI
cd $DESKTOP_DIR/srv
gox -verbose -os="darwin linux" -arch="amd64" -output $WORKING_DIR/bin/{{.OS}}/{{.Arch}}/kazoup-desktop
cd $GOPATH/src/github.com/kazoup/platform/desktop/web
gox -verbose -os="darwin linux" -arch="amd64" -output $WORKING_DIR/bin/{{.OS}}/{{.Arch}}/kazoup-web
cd $GOPATH/src/github.com/micro/micro
gox -verbose -os="darwin linux" -arch="amd64" -output $WORKING_DIR/bin/{{.OS}}/{{.Arch}}/micro
cd $WORKING_DIR


