#!/bin/bash

set -e
set -x
USERNAME=kazoup
REPONAME=platform
WORKING_DIR=$PWD
ghr -t $GITHUB_TOKEN -u $USERNAME -r $REPONAME  --replace v0.0.7 release/
