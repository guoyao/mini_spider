package parser

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"mini_spider/media"
	"mini_spider/util"
)

type WebpageParser struct {
	targetUrl *regexp.Regexp
}

func NewWebpageParser(targetUrl *regexp.Regexp) *WebpageParser {
	return &WebpageParser{targetUrl: targetUrl}
}

func (p *WebpageParser) Parse(m media.Media) ([]*http.Request, error) {
	content := m.Content()
	if content == nil {
		return nil, errors.New("content is nil")
	}

	buf, err := util.ToUTF8(content, m.ContentCharset())
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	links := visit(nil, doc)
	requests := make([]*http.Request, 0, len(links))

	for _, link := range links {
		absUrl, err := util.GetAbsoluteUrl(m.URL(), link)
		if err == nil && absUrl != "" {
			req, err := http.NewRequest("GET", absUrl, nil)
			if err == nil {
				requests = append(requests, req)
			}
		}
	}

	return requests, nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img") {
		for _, a := range n.Attr {
			if a.Key == "href" || a.Key == "src" {
				link := a.Val
				if strings.HasPrefix(link, "javascript:") {
					link = util.GetUrlFromJavascript(link)
				}
				links = append(links, link)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}
