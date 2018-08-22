package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/docopt/docopt-go"

	"mini_spider/config"
	"mini_spider/fetcher"
	"mini_spider/log"
	"mini_spider/spider"
	"mini_spider/storage"
	"mini_spider/util"
)

const VERSION = "0.0.2"

func usage() (map[string]interface{}, error) {
	usage := `
Usage: ./mini_spider [options]

Options:
    -v, --version       	show version
    -h, --help          	show help
    -c CONF_FILE        	set config directory [default: ../conf]
    -l LOG_DIR          	set log directory [default: ../log]
    --depth=<n>			overrite depth config in CONF_FILE
    --storage=<disk|bos>	set storage driver [default: disk]

Example:

    ./mini_spider -c ../conf -l ../log
    ./mini_spider -c /home/xxx/mini_spider/conf -l /home/xxx/mini_spider/log
    ./mini_spider --storage=bos
    ./mini_spider --depth=2`

	return docopt.Parse(usage, nil, true, VERSION, false)
}

func prepare(args map[string]interface{}) *config.Config {
	confDir := args["-c"].(string)
	cfg, err := config.LoadConfigFromFile(confDir + "/spider.conf")
	if err != nil {
		fmt.Println("load config file error: " + err.Error())
		os.Exit(1)
	}

	if args["--depth"] != nil {
		depth := args["--depth"].(string)
		maxDepth, err := strconv.Atoi(depth)
		if err == nil {
			cfg.Spider.MaxDepth = uint(maxDepth)
		}
	}

	logDir := args["-l"].(string)
	err = util.MkdirAll(logDir)
	if err != nil {
		fmt.Println("create log directory error: " + err.Error())
		os.Exit(1)
	}

	log.Init("mini_spider", logDir, "INFO", true)

	return cfg
}

func getStorageDriver(args map[string]interface{}, cfg *config.Config) storage.StorageDriver {
	var driver storage.StorageDriver = storage.NewDiskStorage(cfg.Spider.OutputDirectory)

	if args["--storage"] != nil {
		driverName := args["--storage"].(string)
		if driverName == "bos" {
			bucket := "allitebooks"
			driver = storage.NewBosStorage(os.Getenv("BAIDU_BCE_AK"), os.Getenv("BAIDU_BCE_SK"), bucket)
		}
	}

	return driver
}

func main() {
	args, _ := usage()
	cfg := prepare(args)

	defer log.Logger.Close()

	err := util.MkdirAll(cfg.Spider.OutputDirectory)
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

	targetUrl, err := regexp.Compile(cfg.Spider.TargetUrl)
	if err != nil {
		log.Logger.Critical(err)
		log.Logger.Close()
		os.Exit(1)
	}

	fetcher := fetcher.NewWebpageFetcher(cfg.Spider.CrawlTimeout, getStorageDriver(args, cfg))
	spider := spider.NewSpider(
		targetUrl,
		cfg.Spider.MaxDepth,
		cfg.Spider.CrawlInterval,
		cfg.Spider.ThreadCount,
		fetcher,
	)

	spider.Start(seeds)
}
