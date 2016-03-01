package util

import (
	"log"
	"net/url"
	"os"
)

const (
	// XXX leaving this as http for now.  We probably want to make this unix://
	scheme = "http://"
)

var (
	DefaultHost = Host()
)

func Host() *url.URL {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	thisHost, err := url.Parse(scheme + name)
	if err != nil {
		log.Fatal(err)
	}

	return thisHost
}

// ServiceURL returns the URL for a given service relative to this host
func ServiceURL(serviceName string) *url.URL {
	s, err := DefaultHost.Parse(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
