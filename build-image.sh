#!/bin/bash
VERSION=0.1.0
GOOS=linux GOARCH=amd64 go build -o bin/grpccalc ./server
docker build -t ariannavespri/grpccalc:${VERSION} .
