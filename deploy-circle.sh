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
	if [ $? -ne 0 ]; then
   		 echo "Could not push the image"
      exit 1;
  	fi;
	sudo chown -R ubuntu:ubuntu /home/ubuntu/.kube
	kubectl set image deployment/$IMAGE $IMAGE=eu.gcr.io/desktop-1470249894548/$IMAGE:$CIRCLE_SHA1
	if [ $? -ne 0 ]; then
   		 echo "Could not perform rolling update"
      exit 1;
  	fi;

	popd >/dev/null
done
