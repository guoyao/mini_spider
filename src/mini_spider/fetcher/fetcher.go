/* package fetcher provide different fetchers for the fetch feature of spider */
package fetcher

import (
	"net/http"

	"mini_spider/media"
)

type Fetcher interface {
	GetMetadata(req *http.Request) (media.Media, error)
	Fetch(req *http.Request) (media.Media, error)

	Exist(metadata media.Media) bool
	Save(media media.Media) error
}
