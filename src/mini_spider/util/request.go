package util

import (
	"net/http"
	"regexp"
)

type Request struct {
	*http.Request
	Depth          uint
	ShouldParse    bool
	ShouldDownload bool
}

func NewRequest(req *http.Request, depth uint, shouldParse, shouldDownload bool) *Request {
	return &Request{Request: req, Depth: depth, ShouldParse: shouldParse, ShouldDownload: shouldDownload}
}

func NewRequests(requests []*http.Request, depth uint, targetUrl *regexp.Regexp, hosts []string) []*Request {
	wrappedRequests := make([]*Request, len(requests))

	for index, request := range requests {
		host := request.URL.Hostname()
		shouldParse := false
		for _, value := range hosts {
			if host == value {
				shouldParse = true
				break
			}
		}

		wrappedRequests[index] = NewRequest(request, depth, shouldParse, targetUrl.MatchString(request.URL.String()))
	}

	return wrappedRequests
}
