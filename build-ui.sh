#!/bin/bash

set -e
set -x

CWD="$(pwd)"
# Build UI 
# cleanup
rm -rf  ui/web/html
cd ui/src && npm install && bower install
npm install polymer-cli
node_modules/polymer-cli/bin/polymer.js build -v
cd .. && cp -r src web/html
cd $CWD