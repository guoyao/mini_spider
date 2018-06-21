/* package spider provide the most important structure for spider */
package spider

import (
	"net/http"
	"regexp"
	"sync"
	"time"

	"mini_spider/fetcher"
	"mini_spider/log"
	"mini_spider/parser"
	"mini_spider/util"
)

type Spider struct {
	targetUrl     *regexp.Regexp
	maxDepth      uint
	crawlInterval time.Duration
	threadCount   uint
	fetcher       fetcher.Fetcher
	requestQueue  *util.RequestQueue
}

func NewSpider(targetUrl *regexp.Regexp, maxDepth, interval, threadCount uint, fetcher fetcher.Fetcher) *Spider {
	return &Spider{
		targetUrl:     targetUrl,
		maxDepth:      maxDepth,
		crawlInterval: time.Duration(interval) * time.Second,
		threadCount:   threadCount,
		fetcher:       fetcher,
		requestQueue:  util.NewRequestQueue(),
	}
}

func (s *Spider) addSeeds(seeds []*Seed) {
	for _, seed := range seeds {
		req, err := http.NewRequest("GET", seed.url, nil)
		if err != nil {
			log.Logger.Error("init request for '" + seed.url + "' error: " + err.Error())
			continue
		}
		s.requestQueue.Push(util.NewRequest(req, 0, s.targetUrl.MatchString(seed.url)))
	}
}

func (s *Spider) Start(seeds []*Seed) {
	var wg sync.WaitGroup

	s.addSeeds(seeds)
	routingChan := make(chan struct{}, s.threadCount)

	for {
		if s.requestQueue.Len() == 0 && len(routingChan) == 0 {
			break
		}

		routingChan <- struct{}{}
		wg.Add(1)

		// worker
		go func() {
			defer func() {
				wg.Done()
				<-routingChan
			}()

			req, err := s.requestQueue.Pop()
			if err != nil {
				log.Logger.Error("get request from queue error: " + err.Error())
				return
			}

			if req == nil {
				return
			}

			log.Logger.Info(req.URL.String())

			media, err := s.fetcher.Fetch(req.Request)
			if err != nil {
				log.Logger.Error("fetch error: " + err.Error())
			} else {
				if req.ShouldDownload {
					err = s.fetcher.Save(media)
					if err != nil {
						log.Logger.Error("save to disk error: " + err.Error())
					}
				}

				if req.Depth < s.maxDepth {
					parser := parser.GetParser(media.ContentType(), s.targetUrl)
					if parser != nil {
						requests, err := parser.Parse(media)
						if err != nil {
							log.Logger.Error("parse content error: " + err.Error())
						} else if len(requests) > 0 {
							s.requestQueue.PushAll(util.NewRequests(requests, req.Depth+1, s.targetUrl))
						}
					}
				}
			}

			time.Sleep(s.crawlInterval)
		}()
	}

	wg.Wait()
}
