// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"gopkg.in/urfave/cli.v1"

	"github.com/vmware/vmw-guestinfo/rpcvmx"
	"github.com/vmware/vmw-guestinfo/vmcheck"

	"github.com/vmware/vic/pkg/version"
)

type environment struct {
	Platform   map[string]string
	Properties map[string]string
}

func main() {

	app := cli.NewApp()
	app.Usage = "Fetch OVF environment information"
	app.Version = version.GetBuild().ShortVersion()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "key, k",
			Value: "",
			Usage: "Get single OVF property",
		},
		cli.BoolFlag{
			Name:  "dump",
			Usage: "Dump the OVF Environment XML",
		},
	}

	app.Action = func(c *cli.Context) error {

		// Check if we're running inside a VM
		isVM, err := vmcheck.IsVirtualWorld()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(-1)
		}

		// No point in running if we're not inside a VM
		if !isVM {
			fmt.Println("not living in a virtual world... :(")
			os.Exit(-1)
		}

		config := rpcvmx.NewConfig()
		// Fetch OVF Environment via RPC
		ovfEnv, err := config.String("guestinfo.ovfEnv", "")
		if err != nil {
			fmt.Println("impossible to fetch ovf environment, exiting")
			os.Exit(1)
		}

		if c.Bool("dump") {
			fmt.Println(ovfEnv)
			return nil
		}

		// TODO: fix this when proper support for namespaces is added to golang.
		// ref: golang/go/issues/14407 and golang/go/issues/14407
		ovfEnv = strings.Replace(ovfEnv, "oe:key", "key", -1)
		ovfEnv = strings.Replace(ovfEnv, "oe:value", "value", -1)

		var e environment

		err = xml.Unmarshal([]byte(ovfEnv), &e)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(-1)
		}

		if c.String("key") != "" {
			fmt.Println(e.Properties[c.String("key")])
			return nil
		}

		for k, v := range e.Properties {
			fmt.Printf("[%s]=%s\n", k, v)
		}

		return nil
	}

	app.Run(os.Args)

}

func (e *environment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	type property struct {
		Key   string `xml:"key,attr"`
		Value string `xml:"value,attr"`
	}

	type propertySection struct {
		Property []property `xml:"Property"`
	}
	type platformSection struct {
		Kind    string `xml:"Kind"`
		Version string `xml:"Version"`
		Vendor  string `xml:"Vendor"`
		Locale  string `xml:"Locale"`
	}

	var environment struct {
		XMLName         xml.Name        `xml:"Environment"`
		PlatformSection platformSection `xml:"PlatformSection"`
		PropertySection propertySection `xml:"PropertySection"`
	}
	err := d.DecodeElement(&environment, &start)
	if err == nil {

		e.Platform = make(map[string]string)
		e.Platform["kind"] = environment.PlatformSection.Kind
		e.Platform["version"] = environment.PlatformSection.Version
		e.Platform["vendor"] = environment.PlatformSection.Vendor
		e.Platform["locale"] = environment.PlatformSection.Locale

		e.Properties = make(map[string]string)
		for _, v := range environment.PropertySection.Property {
			e.Properties[v.Key] = v.Value
		}
	}
	return err

}
