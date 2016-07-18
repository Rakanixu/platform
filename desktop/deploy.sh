#!/bin/bash

set -e
set -x
USERNAME=kazoup
REPONAME=platform
WORKING_DIR=$PWD
go get github.com/tcnksm/ghr
ghr -t $GITHUB_TOKEN -u $USERNAME -r $REPONAME --replace `git describe --tags` dist/
