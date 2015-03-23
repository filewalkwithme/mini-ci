#!/bin/bash

mkdir --parents /home/docker/go/src/github.com/$APP
cd /home/docker/go/src/github.com/$APP
git clone https://github.com/$APP.git /home/docker/go/src/github.com/$APP
git checkout -q $COMMIT

if [ $? -eq 0 ]; then
  echo "go build"
  go build
fi

if [ $? -eq 0 ]; then
  echo "go test -v"
  go test -v
fi

if [ $? -eq 0 ]; then
  echo "minideploy"
  minideploy
fi

echo $?
