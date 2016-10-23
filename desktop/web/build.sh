#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
# Build Micro
WORKING_DIR=$GOPATH/src/github.com/kazoup/platform/desktop/web
BUILD_DIR=$GOPATH/src/github.com/kazoup/platform/cmd/kazoup
rm -rf $WORKING_DIR/bin

cd $BUILD_DIR
# crosscompile Kazoup CLI
gox -verbose -os="darwin linux windows" -arch="amd64" -output $WORKING_DIR/bin/{{.OS}}/{{.Arch}}/kazoup
cd $WORKING_DIR


