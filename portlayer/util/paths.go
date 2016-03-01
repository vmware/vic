package util

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	StoragePath = "storage/"
)

// StoreNameToURL parses the image URL in the form /storage/<image store>/<image name>
func StoreNameToURL(storeName string) (*url.URL, error) {
	return ServiceURL(StoragePath).Parse(storeName)
}

func StoreName(u *url.URL) (string, error) {
	// Check the path isn't malformed.
	if !filepath.IsAbs(u.Path) {
		return "", errors.New("invalid uri path")
	}

	segments := strings.Split(filepath.Clean(u.Path), "/")[1:]

	if segments[0] != filepath.Clean(StoragePath) {
		return "", errors.New("not a storage path")
	}

	if len(segments) < 2 {
		return "", errors.New("uri path mismatch")
	}

	return segments[1], nil
}

func ImageURL(storeName, imageName string) (*url.URL, error) {
	u, err := StoreNameToURL(storeName)
	if err != nil {
		return nil, err
	}

	return u.Parse(imageName)
}
