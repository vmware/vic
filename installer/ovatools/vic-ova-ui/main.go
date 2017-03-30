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
//
// +build linux

package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"

	ui "github.com/gizak/termui"

	"github.com/dustin/go-humanize"

	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vmw-guestinfo/rpcvmx"
	"github.com/vmware/vmw-guestinfo/vmcheck"
)

const (
	VT_ACTIVATE   = 0x5606
	VT_WAITACTIVE = 0x5607
)

func main() {

	// If we're running under linux, switch to virtual terminal 2 on startup
	ioctl(uintptr(os.Stdout.Fd()), VT_ACTIVATE, 2)
	ioctl(uintptr(os.Stdout.Fd()), VT_WAITACTIVE, 2)

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	gray := ui.ColorRGB(1, 1, 1)
	blue := ui.ColorCyan

	// Check if we're running inside a VM
	if isVM, err := vmcheck.IsVirtualWorld(); err != nil || !isVM {
		fmt.Fprintln(os.Stderr, "not living in a virtual world... :(")
		os.Exit(-1)
	}

	config := rpcvmx.NewConfig()
	// Fetch OVF Environment via RPC
	ovfEnv, err := config.String("guestinfo.ovfEnv", "")
	if err != nil {
		fmt.Println("impossible to fetch ovf environment, exiting")
		os.Exit(1)
	}

	// TODO: fix this when proper support for namespaces is added to golang.
	// ref: golang/go/issues/14407 and golang/go/issues/14407
	ovfEnv = strings.Replace(ovfEnv, "oe:key", "key", -1)
	ovfEnv = strings.Replace(ovfEnv, "oe:value", "value", -1)

	var ovf environment

	var ips string

	err = xml.Unmarshal([]byte(ovfEnv), &ovf)
	if err != nil {
		ips = fmt.Sprintf("error: %s\n", err.Error())
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// If the IP is a loopback interface, coming from docker0 or an ip6, we don't need it
			if ip.IsLoopback() || ip.String() == "172.17.0.1" || ip.To4() == nil {
				continue
			}
			if strings.ToLower(ovf.Properties["registry.deploy"]) == "true" {
				ips = fmt.Sprintf("%sAccess the Container Registry at:\nhttps://%s:%s\n", ips, ip.String(), ovf.Properties["registry.port"])
			}
			if strings.ToLower(ovf.Properties["management_portal.deploy"]) == "true" {
				ips = fmt.Sprintf("%sAccess the Container Management Portal at:\nhttps://%s:%s\n", ips, ip.String(), ovf.Properties["management_portal.port"])
			}
		}
	}

	bottompanel := ui.NewPar(ips)
	bottompanel.Height = ui.TermHeight()/2 + 1
	bottompanel.Width = ui.TermWidth()
	bottompanel.TextFgColor = ui.ColorBlack
	bottompanel.TextBgColor = blue
	bottompanel.Y = ui.TermHeight() / 2
	bottompanel.X = 0
	bottompanel.Bg = blue
	bottompanel.BorderFg = ui.ColorWhite
	bottompanel.BorderBg = blue
	bottompanel.BorderTop = false
	bottompanel.PaddingTop = 1
	bottompanel.PaddingLeft = 4

	toppanel := ui.NewPar(fmt.Sprintf("VMware vSphere Integrated Containers %s\n\n%s\n%s", version.GetBuild().ShortVersion(), getCPUs(), getMemory()))
	toppanel.Height = ui.TermHeight()/2 + 1
	toppanel.Width = ui.TermWidth()
	toppanel.TextFgColor = ui.ColorWhite
	toppanel.Y = 0
	toppanel.X = 0
	toppanel.TextBgColor = gray
	toppanel.Bg = gray
	toppanel.BorderBg = gray
	toppanel.BorderFg = ui.ColorWhite
	toppanel.BorderBottom = false
	toppanel.PaddingTop = 4
	toppanel.PaddingLeft = 4

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Render(toppanel, bottompanel)

	ui.Loop()
}

type environment struct {
	Properties map[string]string
}

func (e *environment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	type property struct {
		Key   string `xml:"key,attr"`
		Value string `xml:"value,attr"`
	}

	type propertySection struct {
		Property []property `xml:"Property"`
	}

	var environment struct {
		XMLName         xml.Name        `xml:"Environment"`
		PropertySection propertySection `xml:"PropertySection"`
	}
	err := d.DecodeElement(&environment, &start)
	if err == nil {
		e.Properties = make(map[string]string)
		for _, v := range environment.PropertySection.Property {
			e.Properties[v.Key] = v.Value
		}
	}
	return err

}

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}

func getCPUs() string {
	out, _ := exec.Command("lscpu").Output()
	outstring := strings.TrimSpace(string(out))
	lines := strings.Split(outstring, "\n")
	var cpus string
	var model string
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "CPU(s)":
			cpus = value
		case "Model name":
			model = value
		}
	}

	return fmt.Sprintf("%sx %s", cpus, model)
}

func getMemory() string {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	if err != nil {
		panic("Austin, we have a problem... syscall.Sysinfo:" + err.Error())
	}
	return fmt.Sprintf("%s Memory", humanize.IBytes(uint64(si.Totalram)))
}
