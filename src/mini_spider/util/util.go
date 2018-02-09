package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	JS_URL_PATTERN = "javascript:location.href="
)

var charsetRegexp, _ = regexp.Compile("charset=([^\\s]+)")

// URLEncode encodes a string like Javascript's encodeURIComponent()
func URLEncode(str string) string {
	// BUG(go): see https://github.com/golang/go/issues/4013
	// use %20 instead of the + character for encoding a space
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// CheckFileExists checks if specified file exists.
func CheckFileExists(filename string) bool {
	exist := true

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func MkdirAll(path string) error {
	dir, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dir, 0744)
	if err != nil {
		return err
	}

	return nil
}

func GetAbsoluteUrl(baseUrl, relativeUrl string) (string, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	relative, err := url.Parse(relativeUrl)
	if err != nil {
		return "", nil
	}

	return base.ResolveReference(relative).String(), nil
}

func GetCharsetFromContentType(contentType string) string {
	matchGroup := charsetRegexp.FindSubmatch([]byte(contentType))
	if len(matchGroup) > 1 {
		return string(matchGroup[1])
	}

	return ""
}

func ToUTF8(reader io.Reader, encoding string) ([]byte, error) {
	utf8Reader, err := charset.NewReader(reader, encoding)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(utf8Reader)
}

func GetUrlFromJavascript(link string) string {
	if strings.HasPrefix(link, JS_URL_PATTERN) {
		link = strings.Trim(strings.Replace(link, JS_URL_PATTERN, "", 1), "\"'")
	}

	return link
}

func FormatTest(funcName, got, expected string) string {
	return fmt.Sprintf("%s failed. Got %s, expected %s", funcName, got, expected)
}
