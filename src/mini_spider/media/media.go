package media

import (
	"io"
)

type Media interface {
	Name() string
	URL() string
	Content() io.Reader
	ContentType() string
	ContentCharset() string
}
