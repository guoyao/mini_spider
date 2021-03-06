package util

import (
	"net/http"
)

type Request struct {
	*http.Request
	Depth uint
}

func NewRequest(req *http.Request, depth uint) *Request {
	return &Request{Request: req, Depth: depth}
}

func NewRequests(requests []*http.Request, depth uint) []*Request {
	wrappedRequests := make([]*Request, len(requests))

	for index, request := range requests {
		wrappedRequests[index] = NewRequest(request, depth)
	}

	return wrappedRequests
}
