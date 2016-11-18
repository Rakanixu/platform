#!/bin/bash

# Just a script to run the demo
cmd=$1
dir=$2
kube=kubectl
list=(broker registry micro platform )

start() {
	if [ -z $dir ]; then
		for dir in ${list[@]}; do
			find $dir -name "*.yaml" | while read file; do
				$kube create -f $file
			done
		done
		return
	fi

	find $dir -name "*.yaml" | while read file; do
		$kube create -f $file
	done

	# if [ -z $dir ] || [ "$dir" == "database" ]; then
	# 	node=`kubectl get pods | grep pxc-node3 | awk '{print $1}'`
	# 	echo "Run \"kubectl exec $node -i -t -- mysql -u root -p -h pxc-cluster\" and install DBs in database/platform.sql"
	# fi
}

stop() {
	if [ -z $dir ]; then
		for dir in ${list[@]}; do
			find $dir -name "*.yaml" | while read file; do
				$kube delete -f $file
			done
		done
		return
	fi

	find $dir -name "*.yaml" | while read file; do
		$kube delete -f $file
	done	
}

case $cmd in
	start)
	start
	;;
	stop)
	stop
	;;
	restart)
	stop
	start
	;;
	*)
	echo "$0 <start|stop|restart> [dir]"
	exit
	;;
esac