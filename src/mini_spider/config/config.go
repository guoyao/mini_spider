package config

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Spider struct {
		UrlListFile     string // 种子文件路径
		OutputDirectory string // 抓取结果存储目录
		MaxDepth        uint   // 最大抓取深度(种子为0级)
		CrawlInterval   uint   // 抓取间隔. 单位: 秒
		CrawlTimeout    uint   // 抓取超时. 单位: 秒
		TargetUrl       string // 需要存储的目标网页URL pattern(正则表达式)
		ThreadCount     uint   // 抓取routine数
	}
}

func LoadConfigFromFile(filePath string) (*Config, error) {
	var config Config

	err := gcfg.ReadFileInto(&config, filePath)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
