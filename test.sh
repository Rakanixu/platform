#!/bin/bash

set -e
set -x

# Used to test all the things

find * -type d -maxdepth 1 -print | while read dir; do
	if [ ! -f $dir/Dockerfile ]; then
		continue
	fi

	pushd $dir >/dev/null


	SERVICE=${dir%/*}-${dir#*/}
	
	# dep
	go get -d  -v -t ./...
	
	# test
	go test -v ./...


	popd >/dev/null
done

# Frontend tests

# dependencies, add if required
npm install unit.js
npm install mocha
npm install mock-browser
npm install workerjs

mocha ui/polymer/src/**/*-test.js
