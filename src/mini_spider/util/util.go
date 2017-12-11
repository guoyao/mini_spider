package util

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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
