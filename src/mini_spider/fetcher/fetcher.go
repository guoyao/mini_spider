/* package fetcher provide different fetchers for the fetch feature of spider */
package fetcher

import (
	"mini_spider/media"
)

type Fetcher interface {
	Exist(media media.Media) bool
	Save(media media.Media) error
}
