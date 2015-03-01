#!/bin/bash

#TODO: put the git clone stuff here
ls -halt /home/
go version
go get github.com/maiconio/portugo
ls -halt /home/docker/go/src/github.com/maiconio/portugo
cd /home/docker/go/src/github.com/maiconio/portugo
go build
ls -halt /home/docker/go/src/github.com/maiconio/portugo
