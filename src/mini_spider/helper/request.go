package helper

import (
	"net/http"
)

// Request defined a custom request data structure
type Request struct {
	*http.Request
	Depth          uint
	ShouldDownload bool
}

func NewRequest(req *http.Request, depth uint, shouldDownload bool) *Request {
	return &Request{
		Request:        req,
		Depth:          depth,
		ShouldDownload: shouldDownload,
	}
}
