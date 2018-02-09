package util

import (
	"strconv"
	"strings"
	"testing"
)

func TestURLEncode(t *testing.T) {
	expected := "http%3A%2F%2Fwww.baidu.com%2Fexample%2F%E6%B5%8B%E8%AF%95"
	result := URLEncode("http://www.baidu.com/example/测试")

	if result != expected {
		t.Error(FormatTest("URIEncode", result, expected))
	}
}

func TestCheckFileExists(t *testing.T) {
	expected := true
	result := CheckFileExists("util_test.go")

	if result != expected {
		t.Error(FormatTest("CheckFileExists", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	expected = false
	result = CheckFileExists("util_test_2.go")

	if result != expected {
		t.Error(FormatTest("CheckFileExists", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}
}

func TestGetAbsoluteUrl(t *testing.T) {
	expected := "https://www.baidu.com/a/b/c"
	result, err := GetAbsoluteUrl("https://www.baidu.com", "a/b/c")

	if err != nil {
		t.Error("GetAbsoluteUrl failed: " + err.Error())
	} else if result != expected {
		t.Error(FormatTest("GetAbsoluteUrl", result, expected))
	}

	expected = "https://www.baidu.com/a/b/c"
	result, err = GetAbsoluteUrl("https://www.baidu.com/x", "a/b/c")

	if err != nil {
		t.Error("GetAbsoluteUrl failed: " + err.Error())
	} else if result != expected {
		t.Error(FormatTest("GetAbsoluteUrl", result, expected))
	}

	expected = "https://www.baidu.com/a/b/c"
	result, err = GetAbsoluteUrl("https://www.baidu.com/x/y/z", "/a/b/c")

	if err != nil {
		t.Error("GetAbsoluteUrl failed: " + err.Error())
	} else if result != expected {
		t.Error(FormatTest("GetAbsoluteUrl", result, expected))
	}

	expected = "https://www.baidu.com/a/b/c"
	result, err = GetAbsoluteUrl("https://www.baidu.com/x/y", "../a/b/c")

	if err != nil {
		t.Error("GetAbsoluteUrl failed: " + err.Error())
	} else if result != expected {
		t.Error(FormatTest("GetAbsoluteUrl", result, expected))
	}

	expected = "https://www.baidu.com/a/b/c"
	result, err = GetAbsoluteUrl("https://www.baidu.com/x/y", "https://www.baidu.com/a/b/c")

	if err != nil {
		t.Error("GetAbsoluteUrl failed: " + err.Error())
	} else if result != expected {
		t.Error(FormatTest("GetAbsoluteUrl", result, expected))
	}
}

func TestGetCharsetFromContentType(t *testing.T) {
	expected := "UTF-8"
	result := GetCharsetFromContentType("text/html; charset=UTF-8")

	if result != expected {
		t.Error(FormatTest("GetCharsetFromContentType", result, expected))
	}

	expected = ""
	result = GetCharsetFromContentType("text/html")

	if result != expected {
		t.Error(FormatTest("GetCharsetFromContentType", result, expected))
	}

	expected = ""
	result = GetCharsetFromContentType("")

	if result != expected {
		t.Error(FormatTest("GetCharsetFromContentType", result, expected))
	}
}

func TestToUTF8(t *testing.T) {
	expected := "£5 for Peppé"
	utf8Bytes, err := ToUTF8(strings.NewReader("\xa35 for Pepp\xe9"), "latin1")
	result := string(utf8Bytes)

	if err != nil {
		t.Error("ToUTF8 failed: " + err.Error())
	} else if result != expected {
		t.Error(FormatTest("ToUTF8", result, expected))
	}
}

func TestGetUrlFromJavascript(t *testing.T) {
	expected := "page1.html"
	result := GetUrlFromJavascript("javascript:location.href=\"page1.html\"")

	if result != expected {
		t.Error(FormatTest("GetUrlFromJavascript", result, expected))
	}

	expected = "page1.html"
	result = GetUrlFromJavascript("javascript:location.href='page1.html'")

	if result != expected {
		t.Error(FormatTest("GetUrlFromJavascript", result, expected))
	}

	expected = "page1.html"
	result = GetUrlFromJavascript("javascript:location.href=page1.html")

	if result != expected {
		t.Error(FormatTest("GetUrlFromJavascript", result, expected))
	}
}
