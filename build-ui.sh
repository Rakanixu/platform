#!/bin/bash

set -e
set -x


# Build UI 
# cleanup
rm -rf  ui/web/html
cd ui/src && npm install && bower install
npm install polymer-cli
polymer build -v
cd .. && cp -r src/build/bundled web/html