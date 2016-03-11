package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func init() {
	sdk := os.Getenv("GOVC_URL")
	if sdk != "" {
		flag.Set("sdk", sdk)
		flag.Set("vm-path", "docker-appliance")
		flag.Set("cluster", os.Getenv("GOVC_CLUSTER"))
	}

	// fake up a docker-host for pprof collection
	u := url.URL{Scheme: "http", Host: "127.0.0.1:6060"}

	go func() {
		log.Println(http.ListenAndServe(u.Host, nil))
	}()

	flag.Set("docker-host", u.Host)

}

func TestLogTar(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	logFiles = []string{"vicadm.go", "vicadm_test.go"}

	s := &server{
		addr: "127.0.0.1:0",
	}

	err := s.listen()
	if err != nil {
		log.Fatal(err)
	}

	port := s.listenPort()

	go s.serve()

	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/container-logs.tar.gz", port))
	if err != nil {
		t.Fatal(err)
	}

	z, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	tz := tar.NewReader(z)

	for {
		h, err := tz.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
		}

		name, err := url.QueryUnescape(h.Name)
		if err != nil {
			t.Fatal(err)
		}

		if testing.Verbose() {
			fmt.Printf("\n%s...\n", name)
			io.CopyN(os.Stdout, tz, 150)
			fmt.Printf("...\n")
		}
	}
}

func TestLogTail(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
	}

	f, err := ioutil.TempFile("", "vicadm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString("# not much here yet\n")

	logFiles = []string{f.Name(), "vicadm.go"}
	name := filepath.Base(f.Name())

	s := &server{
		addr: "127.0.0.1:0",
	}

	err = s.listen()
	if err != nil {
		log.Fatal(err)
	}

	port := s.listenPort()

	go s.serve()

	out := ioutil.Discard
	if testing.Verbose() {
		out = os.Stdout
	}

	paths := []string{
		"/logs/tail/" + name,
		"/logs/tail",
		"/logs/" + name,
		"/",
	}

	i := 0

	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", port),
	}

	for _, path := range paths {
		u.Path = path
		log.Printf("GET %s:\n", u)
		res, err := http.Get(u.String())
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			for j := 0; j < 512; j++ {
				i++
				f.WriteString(fmt.Sprintf("this is line %d\n", i))
			}
		}()

		size := int64(256)
		n, _ := io.CopyN(out, res.Body, size)
		out.Write([]byte("...\n"))
		res.Body.Close()

		if n != size {
			t.Errorf("expected %d, got %d", size, n)
		}
	}
}
