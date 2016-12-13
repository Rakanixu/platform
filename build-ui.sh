#!/bin/bash

set -e
set -x

CWD="$(pwd)"
# Build UI 
# cleanup
rm -rf  ui/web/html
cd ui/polymer && npm install && bower install
npm install polymer-cli
node_modules/polymer-cli/bin/polymer.js build -v
cd .. && cp -r polymer/build/bundled web/html
cd $CWD