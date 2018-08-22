#!/bin/bash

base_dir=$(cd `dirname $0`; pwd)
export GOPATH="$GOPATH:$base_dir"

go build -o bin/mini_spider mini_spider
