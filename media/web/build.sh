#!/bin/bash

set -e
set -x

# Build Micro
BUILD_DIR=$GOPATH/src/github.com/kazoup/platform/media/web
cd $BUILD_DIR
rm -rf $BUILD_DIR/bin
mkdir -p $BUILD_DIR/bin/darwin/amd64
cd $BUILD_DIR/bin/darwin/amd64 && wget http://evermeet.cx/ffmpeg/ffmpeg-3.1.2.7z && 7z x ffmpeg-3.1.2.7z && rm -rf ffmpeg-3.1.2.7z
cd $BUILD_DIR/bin/darwin/amd64 && wget http://evermeet.cx/ffmpeg/ffprobe-3.1.2.7z && 7z x ffprobe-3.1.2.7z && rm -rf ffprobe-3.1.2.7z

mkdir -p $BUILD_DIR/bin/linux/amd64
cd $BUILD_DIR/bin/linux/amd64 && wget http://johnvansickle.com/ffmpeg/releases/ffmpeg-release-32bit-static.tar.xz && tar -xJf  ffmpeg-release-32bit-static.tar.xz && rm -rf ffmpeg-release-32bit-static.tar.xz

mkdir -p $BUILD_DIR/bin/windows/amd64
cd $BUILD_DIR/bin/windows/amd64 && wget https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-3.0.1-win64-static.7z && 7z x ffmpeg-3.0.1-win64-static.7z && rm -rf ffmpeg-3.0.1-win64-static.7z
cd $BUILD_DIR
