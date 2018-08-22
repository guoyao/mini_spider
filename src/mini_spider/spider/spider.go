/* package spider provide the most important structure for spider */
package spider

import (
	"errors"
	"net/http"
	"regexp"
	"sync"
	"time"

	"mini_spider/fetcher"
	"mini_spider/log"
	"mini_spider/media"
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
		s.requestQueue.Push(util.NewRequest(req, 0, true, s.targetUrl.MatchString(seed.url)))
	}
}

func (s *Spider) Start(seeds []*Seed) {
	var wg sync.WaitGroup

	s.addSeeds(seeds)
	hosts := getHosts(seeds)
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
			if err != nil || req == nil {
				if err != nil {
					log.Logger.Error("get request from queue error: " + err.Error())
				}
				return
			}

			media, err := fetch(s, req)
			if err == nil && media != nil {
				save(s, req, media)
				parse(s, req, media, hosts)
			}

			time.Sleep(s.crawlInterval)
		}()
	}

	wg.Wait()
}

func fetch(s *Spider, req *util.Request) (media.Media, error) {
	if req.ShouldParse || req.ShouldDownload {
		log.Logger.Info(req.URL.String())
		media, err := s.fetcher.Fetch(req.Request)
		if err != nil {
			log.Logger.Error("fetch error: " + err.Error())
		}

		return media, err
	}

	return nil, errors.New("invalid request")
}

func save(s *Spider, req *util.Request, media media.Media) {
	if req.ShouldDownload {
		err := s.fetcher.Save(media)
		if err != nil {
			log.Logger.Error("save error: " + err.Error())
		}
	}
}

func parse(s *Spider, req *util.Request, media media.Media, hosts []string) {
	if req.Depth < s.maxDepth && req.ShouldParse {
		parser := parser.GetParser(media.ContentType(), s.targetUrl)
		if parser != nil {
			requests, err := parser.Parse(media)
			if err != nil {
				log.Logger.Error("parse content error: " + err.Error())
			} else if len(requests) > 0 {
				s.requestQueue.PushAll(util.NewRequests(requests, req.Depth+1, s.targetUrl, hosts))
			}
		}
	}
}
