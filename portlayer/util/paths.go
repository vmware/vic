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

// The image URLs are of the form /storage/<image store>/<image name>
func StoreNameToUrl(storeName string) (*url.URL, error) {
	return ServiceUrl(StoragePath).Parse(storeName)
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

func ImageUrl(storeName, imageName string) (*url.URL, error) {
	u, err := StoreNameToUrl(storeName)
	if err != nil {
		return nil, err
	}

	return u.Parse(imageName)
}
