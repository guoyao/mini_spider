package fetcher

import (
	"net/http"

	"mini_spider/media"
)

type Fetcher interface {
	Fetch(req *http.Request) (media.Media, error)
	Save(media media.Media) error
}
