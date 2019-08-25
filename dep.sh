#!/bin/sh
export GOPATH
GOPATH="`pwd`"
cd src/go-coincheck
dep ensure
dep status
