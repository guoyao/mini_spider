package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"mini_spider/config"
	"mini_spider/fetcher"
	"mini_spider/log"
	"mini_spider/spider"
	"mini_spider/util"
)

const VERSION = "0.0.1"

func usage() (map[string]interface{}, error) {
	usage := `
Usage: ./mini_spider [options]

Options:
    -v, --version       show version
    -h, --help          show help
    -c CONF_FILE        set config directory [default: ../conf]
    -l LOG_DIR          set log directory [default: ../log]

Example:

    ./mini_spider -c ../conf -l ../log
    ./mini_spider -c /home/xxx/mini_spider/conf -l /home/xxx/mini_spider/log`

	return docopt.Parse(usage, nil, true, VERSION, false)
}

func main() {
	args, _ := usage()

	confDir := args["-c"].(string)
	cfg, err := config.LoadConfigFromFile(confDir + "/spider.conf")
	if err != nil {
		fmt.Println("load config file error: " + err.Error())
		os.Exit(1)
	}

	logDir := args["-l"].(string)
	err = util.MkdirAll(logDir)
	if err != nil {
		fmt.Println("create log directory error: " + err.Error())
		os.Exit(1)
	}

	log.Init("mini_spider", logDir, "INFO", true)
	defer log.Logger.Close()

	err = util.MkdirAll(cfg.Spider.OutputDirectory)
	if err != nil {
		log.Logger.Critical(err)
		log.Logger.Close()
		os.Exit(1)
	}

	seeds, err := spider.LoadSeedsFromFile(cfg.Spider.UrlListFile)
	if err != nil {
		log.Logger.Critical(err)
		log.Logger.Close()
		os.Exit(1)
	}

	fetcher := fetcher.NewWebpageFetcher(cfg.Spider.CrawlTimeout, cfg.Spider.OutputDirectory)
	spider := spider.NewSpider(
		cfg.Spider.TargetUrl,
		cfg.Spider.MaxDepth,
		cfg.Spider.CrawlInterval,
		cfg.Spider.ThreadCount,
		fetcher,
	)
	spider.Start(seeds)
}
