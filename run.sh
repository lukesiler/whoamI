#!/bin/sh
PORT=${1}

if [ -z "${PORT}" ]; then
    PORT=8080
fi

docker run --rm -it -p ${PORT}:80 lukesiler/whodat
