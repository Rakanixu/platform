#!/bin/bash

# Build UI

cd ui/frontend && npm install && npm install gulp && bower install && node_modules/gulp/bin/gulp.js && cd ../..


