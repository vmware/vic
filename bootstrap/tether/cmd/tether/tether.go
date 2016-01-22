package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	tetherKey string
	hup       chan os.Signal
)

func init() {
	flag.StringVar(&tetherKey, "key", "/.tether-init/init_key", "tetherd control channel private key")
}

// Hopefully this gets rolled into go.sys at some point
func Mkdev(majorNumber int64, minorNumber int64) int {
	return int((majorNumber << 8) | (minorNumber & 0xff) | ((minorNumber & 0xfff00) << 12))
}

func privateKey(file string) ssh.Signer {
	privateBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to load private key: %v", tetherKey)
		return nil
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", tetherKey)
		return nil
	}
	return private
}

func main() {
	run()
}
