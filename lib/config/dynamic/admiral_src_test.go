package dynamic

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestLookupObjectByTag(t *testing.T) {
	target, err := url.Parse("https://administrator@vsphere.local:Admin!23@10.160.204.51")
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid url %s", os.Args[1])
		os.Exit(1)
	}
	src, _ := NewAdmiralSource(target, true)
	src.Get(nil)
}
