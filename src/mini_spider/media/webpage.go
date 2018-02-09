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
	return &Webpage{name: name, url: url, content: content, contentType: contentType, charset: charset}
}

func (webpage *Webpage) Name() string {
	return webpage.name
}

func (webpage *Webpage) URL() string {
	return webpage.url
}

func (webpage *Webpage) Content() io.Reader {
	if content, ok := webpage.content.(io.Seeker); ok {
		content.Seek(0, io.SeekStart)
	}

	return webpage.content
}

func (webpage *Webpage) ContentType() string {
	return webpage.contentType
}

func (webpage *Webpage) ContentCharset() string {
	return webpage.charset
}
