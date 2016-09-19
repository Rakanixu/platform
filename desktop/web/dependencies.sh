#!/bin/bash

set -e
set -x

go get github.com/mitchellh/gox
WORKING_DIR=$GOPATH/src/github.com/kazoup/platform/desktop/web
BUILD_DIR=$GOPATH/src/github.com/kazoup/platform/desktop/web
rm -rf $BUILD_DIR/elasticsearch
rm -rf $BUILD_DIR/java
# add ES download to build on circle
cd $BUILD_DIR && rm -rf elasticsearch && wget https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip
unzip elasticsearch-2.3.4.zip && rm -rf elasticsearch-2.3.4.zip && mv elasticsearch-2.3.4 elasticsearch
elasticsearch/bin/plugin install mobz/elasticsearch-head
elasticsearch/bin/plugin install delete-by-query
# Download OSX JRE
mkdir -p java/darwin/amd64
wget --no-cookies --no-check-certificate --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u101-b13/jre-8u101-macosx-x64.tar.gz
tar -xzf jre-8u101-macosx-x64.tar.gz
mv jre1.8.0_101.jre/Contents/Home/* java/darwin/amd64
rm -rf jre-8u101-macosx-x64.tar.gz
rm -rf jre1.8.0_101.jre
# Download LINUX JRE
mkdir -p java/linux/amd64
wget --no-cookies --no-check-certificate --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u101-b13/jre-8u101-linux-x64.tar.gz
tar -xzf jre-8u101-linux-x64.tar.gz
mv jre1.8.0_101/* java/linux/amd64
rm -rf jre-8u101-linux-x64.tar.gz
rm -rf jre1.8.0_101
# Download WINDOWS JRE
mkdir -p java/windows/amd64
wget --no-cookies --no-check-certificate --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u101-b13/jre-8u101-windows-x64.tar.gz
tar -xzf jre-8u101-windows-x64.tar.gz
mv jre1.8.0_101/* java/windows/amd64
rm -rf jre-8u101-windows-x64.tar.gz
rm -rf jre1.8.0_101



cd $WORKING_DIR

