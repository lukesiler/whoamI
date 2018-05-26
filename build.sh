#!/bin/sh
#GOOS=linux CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o whodat
docker build -t lukesiler/whodat .
