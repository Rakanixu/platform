#!/bin/bash

set -e
set -x

REGISTRY=kazoup

# Used to rebuild all the things

find * -type d -maxdepth 1 -print | while read dir; do
	if [ ! -f $dir/Dockerfile ]; then
		continue
	fi

	pushd $dir >/dev/null


	IMAGE=${dir%/*}-${dir#*/}

 	
	# push docker image
	sudo /opt/google-cloud-sdk/bin/gcloud docker push eu.gcr.io/${PROJECT_NAME}/$IMAGE
	

	popd >/dev/null
done
