#!/bin/bash

set -e
set -x

REGISTRY=kazoup
go get github.com/mitchellh/gox
# Build Micro
go get github.com/micro/micro
WORKING_DIR=$PWD
# add ES download to build on circle
# https://download.elastic.co/elasticsearch/release/org/elasticsearch/distribution/zip/elasticsearch/2.3.4/elasticsearch-2.3.4.zip

#cd ../../../micro/micro && gox -verbose -os="darwin linux" -arch="386 amd64" -output $WORKING_DIR/bin/micro/{{.OS}}/{{.Arch}}/micro
#cd $WORKING_DIR

#cd ../ui/frontend && npm install && npm install gulp && bower install && node_modules/gulp/bin/gulp.js && cd ../..
# Used to rebuild all the things
cd ../
find * -type d -maxdepth 1 -print | while read dir; do
	if [ ! -f $dir/Dockerfile ]; then
		continue
	fi

	pushd $dir >/dev/null


	IMAGE=${dir%/*}-${dir#*/}
	
	# dep
	go get -d  -v -t ./...
	
	# test
	go test -v ./...

	# build static binary
	#CGO_ENABLED=0 GOOS=linux gox build -a -installsuffix cgo  .
	# crosscompile
	gox -verbose -os="darwin linux" -arch="386 amd64" -output ../../desktop/bin/{{.OS}}/{{.Arch}}/${IMAGE}
	
	# build ui  		
	if [[ ${dir%/*} = "ui" && ${dir#*/} != "static" ]]; then
		mkdir -p $WORKING_DIR/frontend/dist/sections/${dir#*/}
		cp -R -f ../frontend/dist/sections/${dir#*/} $WORKING_DIR/frontend/dist/sections
		IMAGE=${dir#*/}-web

		gox -verbose -os="darwin linux" -arch="386 amd64" -output ../../desktop/bin/${IMAGE}_{{.OS}}_{{.Arch}}
	fi

	if [[ ${dir%/*} = "ui" && ${dir#*/} = "static" ]]; then
		mkdir -p $WORKING_DIR/frontend
		cp -R -f ../frontend/dist $WORKING_DIR/frontend
		IMAGE=${dir#*/}-web
		gox -verbose -os="darwin linux" -arch="386 amd64" -output ../../desktop/bin/${IMAGE}_{{.OS}}_{{.Arch}}
	fi

	# build docker image
	#docker build -t $REGISTRY/$IMAGE .

	
	# remove binary
	#rm ${dir#*/}

	popd >/dev/null
done
