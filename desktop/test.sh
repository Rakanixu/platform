#!/bin/bash

docker run --rm -it -e "APP=auth" -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.6 bash -c '
for GOOS in darwin linux; do
   for GOARCH in 386 amd64; do
     go build -v -o app-$GOOS-$GOARCH
   done
done
'
