package util

import (
	"net/url"
	"strings"
	"testing"
)

func TestServiceUrl(t *testing.T) {
	DefaultHost, _ = url.Parse("http://foo.com/")
	u := ServiceURL(StoragePath)

	if strings.Compare(u.String(), "http://foo.com/storage/") != 0 {
		t.Fail()
		return
	}
}
