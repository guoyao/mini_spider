#!/bin/bash

base_dir=$(cd `dirname $0`; pwd)
export GOPATH="$GOPATH:$base_dir"

cd src
go run mini_spider.go
