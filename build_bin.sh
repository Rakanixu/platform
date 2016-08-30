#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
# Build Micro
WORKING_DIR=$GOPATH/src/github.com/kazoup/platform
BUILD_DIR=$GOPATH/src/github.com/kazoup/platform/build
rm -rf $BUILD_DIR
mkdir $BUILD_DIR
# add ES download to build on circle
cd $BUILD_DIR && rm -rf elasticsearch && wget https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip 
unzip elasticsearch-2.3.4.zip && rm -rf elasticsearch-2.3.4.zip && mv elasticsearch-2.3.4 elasticsearch
elasticsearch/bin/plugin install mobz/elasticsearch-head
cd $WORKING_DIR/cmd/kazoup
# crosscompile Kazoup CLI
gox -verbose -os="darwin linux windows" -arch="amd64" -output $BUILD_DIR/bin/{{.OS}}/{{.Arch}}/kazoup
cd $WORKING_DIR

cd $BUILD_DIR/bin/darwin/amd64 && wget http://evermeet.cx/ffmpeg/ffmpeg-3.1.2.7z && 7z x ffmpeg-3.1.2.7z && rm -rf ffmpeg-3.1.2.7z
cd $BUILD_DIR/bin/darwin/amd64 && wget http://evermeet.cx/ffmpeg/ffprobe-3.1.2.7z && 7z x ffprobe-3.1.2.7z && rm -rf ffprobe-3.1.2.7z
cd $BUILD_DIR/bin/linux/amd64 && wget http://johnvansickle.com/ffmpeg/releases/ffmpeg-release-32bit-static.tar.xz && tar -xJf  ffmpeg-release-32bit-static.tar.xz && rm -rf ffmpeg-release-32bit-static.tar.xz
cd $BUILD_DIR/bin/windows/amd64 && wget https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-3.0.1-win64-static.7z && 7z x ffmpeg-3.0.1-win64-static.7z && rm -rf ffmpeg-3.0.1-win64-static.7z

