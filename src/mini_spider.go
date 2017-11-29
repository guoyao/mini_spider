package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"config"
)

const VERSION = "0.0.1"

func usage() (map[string]interface{}, error) {
	usage := `
Usage: ./mini_spider [options]

Options:
    -v, --version       show version
    -h, --help          show help
    -c CONF_FILE        set config directory [default: ./conf]
    -l LOG_DIR          set log directory [default: ./log]

Example:

    ./mini_spider -c ./spider.conf -l ./log
    ./mini_spider -c /home/xxx/mini_spider/spider.conf -l /home/xxx/mini_spider/log`

	return docopt.Parse(usage, nil, true, VERSION, false)
}

func main() {
	args, err := usage()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	confDir := args["-c"].(string)
	config, err := config.LoadConfigFromFile(confDir + "/spider.conf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println((*config).Spider.TargetUrl)
}
