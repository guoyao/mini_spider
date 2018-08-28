package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"

	"mini_spider/config"
	"mini_spider/fetcher"
	"mini_spider/log"
	"mini_spider/spider"
	"mini_spider/storage"
	"mini_spider/util"
)

const VERSION = "0.0.3"

func usage() (map[string]interface{}, error) {
	usage := `
Usage: ./mini_spider [options]

Options:
    -v, --version       	show version
    -h, --help          	show help
    -c CONF_FILE        	set config directory [default: ../conf]
    -l LOG_DIR          	set log directory [default: ../log]
    --output=<outputDirectory>	overrite outputDirectory config in CONF_FILE
    --depth=<maxDepth>		overrite depth config in CONF_FILE
    --interval=<crawlInterval>	overrite crawlInterval config in CONF_FILE
    --timeout=<crawlTimeout>	overrite crawlTimeout config in CONF_FILE
    --target-url=<targetUrl>	overrite targetUrl config in CONF_FILE
    --thread=<threadCount>	overrite threadCount config in CONF_FILE
    --storage=<disk|bos>	set storage driver [default: disk]
    --seed=<seed>		override seed config in data/url.data

Example:

    ./mini_spider -v
    ./mini_spider -h
    ./mini_spider -c ../conf -l ../log
    ./mini_spider -c /home/xxx/mini_spider/conf -l /home/xxx/mini_spider/log
    ./mini_spider --output=./output --depth=2 --interval=2 --timeout=60 --thread=10 --storage=bos
    ./mini_spider --target-url=(java|spring|go).*\\.(pdf)(\\?.+)?$
    ./mini_spider --seed=http://www.baidu.com,http://www.google.com`

	return docopt.Parse(usage, nil, true, VERSION, false)
}

func prepare(args map[string]interface{}) *config.Config {
	confDir := args["-c"].(string)
	cfg, err := config.LoadConfigFromFile(confDir + "/spider.conf")
	if err != nil {
		fmt.Println("load config file error: " + err.Error())
		os.Exit(1)
	}

	if args["--output"] != nil {
		cfg.Spider.OutputDirectory = args["--output"].(string)
	}

	if args["--depth"] != nil {
		depth := args["--depth"].(string)
		maxDepth, err := strconv.Atoi(depth)
		if err == nil {
			cfg.Spider.MaxDepth = uint(maxDepth)
		}
	}

	if args["--interval"] != nil {
		interval := args["--interval"].(string)
		crawlInterval, err := strconv.Atoi(interval)
		if err == nil {
			cfg.Spider.CrawlInterval = uint(crawlInterval)
		}
	}

	if args["--timeout"] != nil {
		timeout := args["--timeout"].(string)
		crawlTimeout, err := strconv.Atoi(timeout)
		if err == nil {
			cfg.Spider.CrawlTimeout = uint(crawlTimeout)
		}
	}

	if args["--target-url"] != nil {
		cfg.Spider.TargetUrl = args["--target-url"].(string)
	}

	if args["--thread"] != nil {
		thread := args["--thread"].(string)
		threadCount, err := strconv.Atoi(thread)
		if err == nil {
			cfg.Spider.ThreadCount = uint(threadCount)
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

func getSeeds(args map[string]interface{}, cfg *config.Config) []*spider.Seed {
	if args["--seed"] != nil {
		seed := strings.TrimSpace(args["--seed"].(string))
		if len(seed) > 0 {
			urls := strings.Split(seed, ",")
			if length := len(urls); length > 0 {
				seeds := make([]*spider.Seed, 0, length)
				for _, v := range urls {
					if url := strings.TrimSpace(v); url != "" {
						seeds = append(seeds, spider.NewSeed(url))
					}
				}
				if len(seeds) > 0 {
					return seeds
				}
			}
		}
	}

	seeds, err := spider.LoadSeedsFromFile(cfg.Spider.UrlListFile)
	if err != nil {
		log.Logger.Critical(err)
		log.Logger.Close()
		os.Exit(1)
	}

	return seeds
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

	seeds := getSeeds(args, cfg)
	if len(seeds) == 0 {
		log.Logger.Critical("no seeds specified")
		log.Logger.Close()
		os.Exit(1)
	}

	targetUrl, err := regexp.Compile(cfg.Spider.TargetUrl)
	if err != nil {
		log.Logger.Critical(err)
		log.Logger.Close()
		os.Exit(1)
	}

	fetcher := fetcher.NewHttpFetcher(cfg.Spider.CrawlTimeout, getStorageDriver(args, cfg))
	spider := spider.NewSpider(
		targetUrl,
		cfg.Spider.MaxDepth,
		cfg.Spider.CrawlInterval,
		cfg.Spider.ThreadCount,
		fetcher,
	)

	spider.Start(seeds)
}
