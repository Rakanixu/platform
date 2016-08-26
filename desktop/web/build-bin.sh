#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
# Build Micro
go get github.com/micro/micro
WORKING_DIR=$PWD
DESKTOP_DIR=$GOPATH/src/github.com/kazoup/platform/desktop
# add ES download to build on circle
cd $DESKTOP_DIR/web && rm -rf elasticsearch && wget https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip 
unzip elasticsearch-2.3.4.zip && rm -rf elasticsearch-2.3.4.zip && mv elasticsearch-2.3.4 elasticsearch
elasticsearch/bin/plugin install mobz/elasticsearch-head

#cd ../../../micro/micro && gox -verbose -os="darwin linux" -arch="386 amd64" -output $WORKING_DIR/bin/micro/{{.OS}}/{{.Arch}}/micro
cd $WORKING_DIR

#cd ../ui/frontend && npm install && npm install gulp && bower install && node_modules/gulp/bin/gulp.js && cd ../..

# Build binary distributable
rm -rf  $GOPATH/src/github.com/kazoup/platform/desktop/web/bin
# crosscompile Kazoup CLI
cd $DESKTOP_DIR/srv
gox -verbose -os="darwin linux windows" -arch="amd64" -output $DESKTOP_DIR/web/bin/{{.OS}}/{{.Arch}}/kazoup-desktop
cd $GOPATH/src/github.com/kazoup/platform/desktop/web
gox -verbose -os="darwin linux windows" -arch="amd64" -output $DESKTOP_DIR/web/bin/{{.OS}}/{{.Arch}}/kazoup-web
cd $GOPATH/src/github.com/micro/micro
gox -verbose -os="darwin linux windows" -arch="amd64" -output $DESKTOP_DIR/web/bin/{{.OS}}/{{.Arch}}/micro
cd $WORKING_DIR
cd $DESKTOP_DIR/web/bin/darwin/amd64 && wget http://evermeet.cx/ffmpeg/ffmpeg-3.1.2.7z && 7z x ffmpeg-3.1.2.7z && rm -rf ffmpeg-3.1.2.7z
cd $DESKTOP_DIR/web/bin/darwin/amd64 && wget http://evermeet.cx/ffmpeg/ffprobe-3.1.2.7z && 7z x ffprobe-3.1.2.7z && rm -rf ffprobe-3.1.2.7z
cd $DESKTOP_DIR/web/bin/linux/amd64 && wget http://johnvansickle.com/ffmpeg/releases/ffmpeg-release-32bit-static.tar.xz && tar -xJf  ffmpeg-release-32bit-static.tar.xz && rm -rf ffmpeg-release-32bit-static.tar.xz
cd $DESKTOP_DIR/web/bin/windows/amd64 && wget https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-3.0.1-win64-static.7z && 7z x ffmpeg-3.0.1-win64-static.7z && rm -rf ffmpeg-3.0.1-win64-static.7z

#cd $WORKING_DIR
