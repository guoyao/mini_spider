package spider

import (
	"net/http"
	"sync"
	"time"

	"mini_spider/fetcher"
	"mini_spider/log"
	"mini_spider/util"
)

type Spider struct {
	maxDepth      uint
	crawlInterval time.Duration
	threadCount   uint
	seeds         []*Seed
	fetcher       fetcher.Fetcher
	requestQueue  *util.RequestQueue
}

func NewSpider(maxDepth, interval, threadCount uint, seeds []*Seed, fetcher fetcher.Fetcher) *Spider {
	return &Spider{
		maxDepth:      maxDepth,
		crawlInterval: time.Duration(interval) * time.Second,
		threadCount:   threadCount,
		seeds:         seeds,
		fetcher:       fetcher,
		requestQueue:  util.NewRequestQueue(),
	}
}

func (spider *Spider) Start() {
	for _, seed := range spider.seeds {
		req, err := http.NewRequest("GET", seed.url, nil)
		if err != nil {
			log.Logger.Error("init request for '" + seed.url + "' error: " + err.Error())
			continue
		}
		spider.requestQueue.Push(req)
	}

	routingChan := make(chan struct{}, spider.threadCount)
	var wg sync.WaitGroup
	for {
		if spider.requestQueue.Len() == 0 {
			break
		}

		routingChan <- struct{}{}
		wg.Add(1)

		go func() {
			defer func() {
				wg.Done()
				<-routingChan
			}()

			req, err := spider.requestQueue.Pop()
			if err != nil {
				log.Logger.Error("get request from queue error: " + err.Error())
				return
			}

			if req == nil {
				return
			}

			media, err := spider.fetcher.Fetch(req)
			if err != nil {
				log.Logger.Error("fetch error: " + err.Error())
			} else {
				err = spider.fetcher.Save(media)
				if err != nil {
					log.Logger.Error("save to disk error: " + err.Error())
				}
			}
			time.Sleep(spider.crawlInterval)
		}()
	}
	wg.Wait()
}
