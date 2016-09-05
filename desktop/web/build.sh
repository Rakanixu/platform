#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
# Build Micro
WORKING_DIR=$GOPATH/src/github.com/kazoup/platform/desktop/web
BUILD_DIR=$GOPATH/src/github.com/kazoup/platform/desktop/web
BIN_DIR=$GOPATH/src/github.com/kazoup/platform/cmd/kazoup
rm -rf $BUILD_DIR/bin
rm -rf $BUILD_DIR/elasticsearch
# add ES download to build on circle
cd $BUILD_DIR && rm -rf elasticsearch && wget https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip
unzip elasticsearch-2.3.4.zip && rm -rf elasticsearch-2.3.4.zip && mv elasticsearch-2.3.4 elasticsearch
elasticsearch/bin/plugin install mobz/elasticsearch-head
cd $BIN_DIR
# crosscompile Kazoup CLI
gox -verbose -os="darwin linux windows" -arch="amd64" -output $BUILD_DIR/bin/{{.OS}}/{{.Arch}}/kazoup
cd $WORKING_DIR


