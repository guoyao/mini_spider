/* package media provide different type contents for the fetch resource of spider */
package media

import (
	"io"
)

type Media interface {
	Name() string
	URL() string
	Content() io.Reader
	SetContent(io.Reader)
	ContentType() string
	ContentLength() uint
	ContentCharset() string
}
