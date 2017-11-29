package main

import (
	"fmt"
	"os"
)

const (
	VERSION      = "0.0.1"
	usageContent = `
Usage: 

    ./mini_spider [options]

    Options:

    -c="./spider.conf"        Set the path of configuration file
    -h                        Show help
    -l="./log"                Set the directory of log 
    -v                        Print version

Example:

    ./mini_spider -c ./spider.conf -l ./log
    ./mini_spider -c /home/xxx/mini_spider/spider.conf -l /home/xxx/mini_spider/log 
`
)

func usage() {
	fmt.Println(usageContent)
	os.Exit(0)
}

func main() {
	usage()
}
