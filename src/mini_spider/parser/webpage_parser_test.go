package parser

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"mini_spider/media"
	"mini_spider/util"
)

func TestParse(t *testing.T) {
	funcName := "Parse"
	html := `
<!DOCTYPE html>
<html>
    <head>
        <meta charset=utf8>
        <title>Crawl Me</title>
    </head>
    <body>
        <ul>
            <li><a href=page1.html>page 1</a></li>
            <li><a href="page2.html">page 2</a></li>
            <li><a href='page3.html'>page 3</a></li>
            <li><a href='mirror/index.html'>mirror</a></li>
            <li><a href='javascript:location.href="page4.html"'>page 4</a></li>
        </ul>
    </body>
</html>`

	targetUrl, err := regexp.Compile(".*.(htm|html)$")
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	parser := NewWebpageParser(targetUrl)
	webpage := media.NewWebpage("", strings.NewReader(html), "", "")
	requests, err := parser.Parse(webpage)
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	expected := 5
	result := len(requests)
	if result != expected {
		t.Error(util.FormatTest(funcName, strconv.Itoa(result), strconv.Itoa(expected)))
	}
}
