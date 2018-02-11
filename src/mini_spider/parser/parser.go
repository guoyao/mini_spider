/* package parser provide different content parser for the parse feature of spider */
package parser

import (
	"net/http"
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

func GetParserByContentType(contentType, urlPattern string) (Parser, error) {
	switch {
	case strings.HasPrefix(contentType, HTML):
		if parser, ok := parserMap[HTML]; ok {
			return parser, nil
		}

		parser, err := NewWebpageParser(urlPattern)
		if err != nil {
			return nil, err
		}

		parserMap[HTML] = parser
		return parser, nil
	}

	return nil, nil
}
