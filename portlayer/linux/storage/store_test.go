package storage

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/stretchr/testify/assert"
	portlayer "github.com/vmware/vic/portlayer/storage"
)

func TestGetImageStoreMissing(t *testing.T) {
	s := &portlayer.NameLookupCache{
		DataStore: &LocalStore{},
	}

	u, err := s.GetImageStore("doesntexist")
	if !assert.Error(t, err, "An error was expected") {
		return
	}

	if !assert.Nil(t, u) {
		return
	}
}

func TestCreateAndGetStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "parentStore")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(dir)

	ls := &LocalStore{
		Path: dir,
	}

	s := &portlayer.NameLookupCache{
		DataStore: ls,
	}

	// test the create through the cache
	createURL, err := s.CreateImageStore("testStore")
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, createURL) {
		return
	}

	// test the store got created on the local store
	localURL, err := ls.GetImageStore("testStore")
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, localURL) {
		return
	}

	if !assert.Equal(t, createURL.String(), localURL.String(), "urls should be equal") {
		return
	}

	// test the cache gives us the same url
	cacheURL, err := s.GetImageStore("testStore")
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, cacheURL) {
		return
	}

	if !assert.Equal(t, createURL.String(), cacheURL.String(), "urls should be equal") {
		return
	}
}

// Creates a tar archive in memory and uses this to test image creation
func TestCreateImage(t *testing.T) {
	dir, err := ioutil.TempDir("", "parentStore")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(dir)

	ls := &LocalStore{
		Path: dir,
	}

	s := &portlayer.NameLookupCache{
		DataStore: ls,
	}

	storeURL, err := s.CreateImageStore("testStore")
	if !assert.NoError(t, err) {
		return
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name string
		Type byte
		Body string
	}{
		{"dir1", tar.TypeDir, ""},
		{"dir1/readme.txt", tar.TypeReg, "This archive contains some text files."},
		{"dir1/gopher.txt", tar.TypeReg, "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"dir1/todo.txt", tar.TypeReg, "Get animal handling license."},
	}

	for _, file := range files {
		hdr := &tar.Header{
			Name:     file.Name,
			Mode:     0777,
			Typeflag: file.Type,
			Size:     int64(len(file.Body)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}

		if file.Type == tar.TypeDir {
			continue
		}

		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}

	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}

	// base this image off scratch
	scratch, err := s.GetImage(storeURL, portlayer.Scratch.ID)
	if !assert.NoError(t, err) {
		return
	}

	// XXX HACK
	tarOptions = &archive.TarOptions{NoLchown: true}

	newImage, err := ls.WriteImage(scratch, "EEEEEE", buf)
	if !assert.NoError(t, err) || !assert.NotNil(t, newImage) {
		return
	}

	// verify we did anything
}
