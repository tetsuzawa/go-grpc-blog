#!/usr/bin/env bash

PROTO_DIR=${HOME}/dev/protocols
PROTO="$1"

#echo "${PROTO}"

mkdir -p ./protocols/"${PROTO}"
protoc -I "${PROTO_DIR}"/"${PROTO}"/ "${PROTO_DIR}"/"${PROTO}"/"${PROTO}".proto --go_out=plugins=grpc:./protocols/"${PROTO}"/
