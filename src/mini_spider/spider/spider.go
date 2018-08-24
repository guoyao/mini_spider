/* package spider provide the most important structure for spider */
package spider

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"mini_spider/fetcher"
	"mini_spider/helper"
	"mini_spider/log"
	"mini_spider/media"
	"mini_spider/parser"
)

// Spider defined a spider data structure
type Spider struct {
	targetUrl     *regexp.Regexp
	maxDepth      uint
	crawlInterval time.Duration
	threadCount   uint
	fetcher       fetcher.Fetcher
	requestQueue  *helper.RequestQueue
}

// NewSpider return a spider instance
func NewSpider(targetUrl *regexp.Regexp, maxDepth, interval, threadCount uint, fetcher fetcher.Fetcher) *Spider {
	return &Spider{
		targetUrl:     targetUrl,
		maxDepth:      maxDepth,
		crawlInterval: time.Duration(interval) * time.Second,
		threadCount:   threadCount,
		fetcher:       fetcher,
		requestQueue:  helper.NewRequestQueue(),
	}
}

// Start crawl
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

			media, err := fetch(s, req, hosts)
			if err == nil && media != nil {
				save(s, req, media)
				parse(s, req, media)
			}

			time.Sleep(s.crawlInterval)
		}()
	}

	wg.Wait()
}

func (s *Spider) addSeeds(seeds []*Seed) {
	for _, seed := range seeds {
		req, err := http.NewRequest("GET", seed.url, nil)
		if err != nil {
			log.Logger.Error("init request for '" + seed.url + "' error: " + err.Error())
			continue
		}
		s.requestQueue.Push(helper.NewRequest(req, 0, s.targetUrl.MatchString(strings.ToLower(seed.url))))
	}
}

func fetch(s *Spider, req *helper.Request, hosts []string) (media.Media, error) {
	log.Logger.Info(req.URL.String())

	needBody := func(media media.Media) bool {
		shouldParse := getShouldParse(req.Request.URL.Hostname(), media.ContentType(), hosts, s.targetUrl)
		shouldDownload := s.targetUrl.MatchString(strings.ToLower(req.Request.URL.String()))
		exists := s.fetcher.Exist(media)
		if exists {
			log.Logger.Info("skip exist: " + req.URL.String())
		}
		return !exists && (shouldParse || shouldDownload)
	}

	if httpFetcher, ok := s.fetcher.(*fetcher.HttpFetcher); ok {
		media, err := httpFetcher.Fetch(req.Request, needBody)
		if err != nil {
			log.Logger.Error("fetch error: " + err.Error())
		}
		return media, err
	}

	return nil, errors.New("no fetcher")
}

func save(s *Spider, req *helper.Request, media media.Media) {
	if req.ShouldDownload && !s.fetcher.Exist(media) {
		err := s.fetcher.Save(media)
		if err != nil {
			log.Logger.Error("save error: " + err.Error())
		}
	}
}

func parse(s *Spider, req *helper.Request, media media.Media) {
	if req.Depth < s.maxDepth {
		parser := parser.GetParser(media.ContentType(), s.targetUrl)
		if parser != nil {
			requests, err := parser.Parse(media)
			if err != nil {
				log.Logger.Error("parse content error: " + err.Error())
			} else if len(requests) > 0 {
				requestSlice := make([]*helper.Request, len(requests))

				for i, request := range requests {
					shouldDownload := s.targetUrl.MatchString(
						strings.ToLower(request.URL.String()))
					requestSlice[i] = helper.NewRequest(request, req.Depth+1, shouldDownload)
				}

				s.requestQueue.PushAll(requestSlice)
			}
		}
	}
}

func getShouldParse(host, contentType string, hosts []string, targetUrl *regexp.Regexp) bool {
	isValidHost := false
	for _, value := range hosts {
		if host == value {
			isValidHost = true
			break
		}
	}

	parser := parser.GetParser(contentType, targetUrl)
	return parser != nil && isValidHost
}
