#!/usr/bin/env bash

BUILDDATE=$(date +%Y%m%d)

GOOS=linux ARCH=amd64 go build -ldflags "-s -w" -v -o bin/wolfcli-${BUILDDATE}