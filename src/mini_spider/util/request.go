package util

import (
	"net/http"
	"regexp"
)

type Request struct {
	*http.Request
	Depth          uint
	ShouldDownload bool
}

func NewRequest(req *http.Request, depth uint, shouldDownload bool) *Request {
	return &Request{Request: req, Depth: depth, ShouldDownload: shouldDownload}
}

func NewRequests(requests []*http.Request, depth uint, targetUrl *regexp.Regexp) []*Request {
	wrappedRequests := make([]*Request, len(requests))

	for index, request := range requests {
		wrappedRequests[index] = NewRequest(request, depth, targetUrl.MatchString(request.URL.String()))
	}

	return wrappedRequests
}
