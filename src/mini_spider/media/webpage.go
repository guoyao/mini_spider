package media

import (
	"io"

	"mini_spider/util"
)

type Webpage struct {
	name        string
	url         string
	content     io.Reader
	contentType string
	charset     string
}

func NewWebpage(url string, content io.Reader, contentType, charset string) *Webpage {
	name := util.URLEncode(url)
	return &Webpage{name: name, url: url, content: content, contentType: contentType,
		charset: charset}
}

func (w *Webpage) Name() string {
	return w.name
}

func (w *Webpage) URL() string {
	return w.url
}

func (w *Webpage) Content() io.Reader {
	if content, ok := w.content.(io.Seeker); ok {
		content.Seek(0, io.SeekStart)
	}

	return w.content
}

func (w *Webpage) ContentType() string {
	return w.contentType
}

func (w *Webpage) ContentCharset() string {
	return w.charset
}
