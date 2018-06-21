/* package parser provide different content parser for the parse feature of spider */
package parser

import (
	"net/http"
	"regexp"
	"strings"

	"mini_spider/media"
)

type Parser interface {
	Parse(media media.Media) ([]*http.Request, error)
}

const (
	HTML = "text/html"
)

var parserMap map[string]Parser = make(map[string]Parser)

func GetParser(contentType string, targetUrl *regexp.Regexp) Parser {
	switch {
	case strings.HasPrefix(contentType, HTML):
		if parser, ok := parserMap[HTML]; ok {
			return parser
		}

		parser := NewWebpageParser(targetUrl)
		parserMap[HTML] = parser
		return parser
	}

	return nil
}
