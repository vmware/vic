// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	priv, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", tetherKey)
		return nil
	}
	return priv
}

func main() {
	run()
}
