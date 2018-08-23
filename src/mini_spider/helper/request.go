package helper

import (
	"net/http"

	"mini_spider/media"
)

// Request defined a custom request data structure
type Request struct {
	*http.Request
	Depth          uint
	ShouldParse    bool
	ShouldDownload bool
	Metadata       media.Media
}

func NewRequest(req *http.Request, depth uint, shouldParse, shouldDownload bool, metadata media.Media) *Request {
	return &Request{
		Request:        req,
		Depth:          depth,
		ShouldParse:    shouldParse,
		ShouldDownload: shouldDownload,
		Metadata:       metadata,
	}
}
