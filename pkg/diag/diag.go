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

package diag

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vmware/vic/pkg/trace"
)

const StatusCodeFatalThreshold = 64

const (
	// Host name was resolved and all pings went through.
	PingStatusOk = 0
	// Host name was resolved, some of the ping responses didn't go through.
	PingStatusOkPacketLosses = 1
	// Host name was resolved, ping didn't return any responses.
	PingStatusOkNotPingable = 2
	// Host name is unknown.
	PingStatusResolutionFailed = 64
	// ping command failed to run.
	PingStatusPingNotExists = 65
	// Ping didn't return anything
	PingStatusNoPingOutput = 66
	// Return in cases where there is no expected output from ping command.
	PingStatusUnknownError = 67
)

const (
	// VC/ESXi API is available.
	VCStatusOK = 0
	// Provided VC/ESXi API URL is wrong.
	VCStatusInvalidURL = 64
	// Error happened trying to query VC/ESXi API
	VCStatusErrorQuery = 65
	// Received response doesn't contain expected data.
	VCStatusErrorResponse = 66
	// Received in case if returned data from server is different from expected.
	VCStatusIncorrectResponse = 67
	// Received response is not XML
	VCStatusNotXML = 68
)

func UserReadablePingTestDescription(code int) string {
	switch code {
	case PingStatusOk:
		return "is reacheable and responds to ICMP ping request"
	case PingStatusOkPacketLosses:
		return "is reacheable and responds to ICMP ping request, some of requests were not answered"
	case PingStatusOkNotPingable:
		return "was resolved, no responses were received to ICMP ping request"
	case PingStatusResolutionFailed:
		return "is unknown host name"
	case PingStatusPingNotExists:
		return "ping command is not available on VCH. Broken image."
	case PingStatusNoPingOutput:
		return "could not be tested due to no ping command output"
	case PingStatusUnknownError:
		return "could not be tested due to unknown error"
	default:
		return "could not be tested due to unknown error code. It may happen if you are using obsolete version of the client."
	}
}

func UserReadableVCAPITestDescription(code int) string {
	switch code {
	case VCStatusOK:
		return "responds as expected"
	case VCStatusInvalidURL:
		return "url is invalid"
	case VCStatusErrorQuery:
		return "failed to respond to the query"
	case VCStatusIncorrectResponse:
		return "returns unexpected response"
	case VCStatusErrorResponse:
		return "returns error"
	case VCStatusNotXML:
		return "returns non XML response"
	default:
		return "unknown API test"
	}
}

var statsRE = *regexp.MustCompile("\\d+ \\w{6,8} transmitted, (\\d+) ")

// CheckPing runs ping to check if VC/ESXi how is resolvable and pingable.
func CheckPing(hn string) int {
	cmd := exec.Command("ping", "-c", "4", "-W", "3", "-i", "0.1", hn)
	return runPing(cmd.CombinedOutput())
}

func runPing(data []byte, err error) int {
	op := trace.NewOperation(context.Background(), "ping test")
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") ||
			strings.Contains(err.Error(), "no such file or directory") ||
			strings.Contains(err.Error(), "system cannot find the file specified") {
			op.Errorf("Ping command not found")
			return PingStatusPingNotExists
		}
	}

	text := string(data)
	if text == "" {
		op.Errorf("Ping didn't return any output")
		return PingStatusNoPingOutput
	}

	if strings.Contains(strings.ToLower(text), "unknown host") ||
		strings.Contains(text, "cannot resolve") {
		op.Errorf("Ping didn't resolve host name")
		return PingStatusResolutionFailed
	}

	found := statsRE.FindAllStringSubmatch(text, -1)
	if len(found) != 1 || len(found[0]) != 2 {
		// this is very unlikely. However, if it happens it is better to know details.
		op.Errorf("Ping returned unexpected response: %s", text)
		return PingStatusUnknownError
	}
	var receiveCnt int64
	if receiveCnt, err = strconv.ParseInt(found[0][1], 10, 64); err != nil {
		op.Errorf("Unknown ping output: %s", err)
		return PingStatusUnknownError
	}

	if receiveCnt > 0 && receiveCnt < 4 {
		op.Errorf("Potential packet losses. Received %d instead of 4 responses.", receiveCnt)
		return PingStatusOkPacketLosses
	}

	if receiveCnt == 0 {
		op.Errorf("ICMP protocol might be blocked. No ping resonses were received.")
		return PingStatusOkNotPingable
	}
	op.Infof("Ping was executed succesfully")
	return PingStatusOk
}

// CheckAPIAvailability accesses VC/ESXi API to ensure it is a correct end point that is up and running.
func CheckAPIAvailability(targetURL string) int {
	op := trace.NewOperation(context.Background(), "api test")
	errorCode := VCStatusErrorQuery

	u, err := url.Parse(targetURL)
	if err != nil {
		return VCStatusInvalidURL
	}

	u.Path = "/sdk/vimService.wsdl"
	apiURL := u.String()

	op.Debugf("Checking access to: %s", apiURL)

	for attempts := 5; errorCode != VCStatusOK && attempts > 0; attempts-- {

		c := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			// Is 20 seconds enough to receive any response from ESXi server?
			Timeout: time.Second * 20,
		}
		errorCode = queryAPI(op, c.Get, apiURL)
	}
	return errorCode
}

func queryAPI(op trace.Operation, getter func(string) (*http.Response, error), apiURL string) int {
	resp, err := getter(apiURL)
	if err != nil {
		op.Errorf("Query error: %s", err)
		return VCStatusErrorQuery
	}

	data := make([]byte, 65636)
	n, err := io.ReadFull(resp.Body, data)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		op.Errorf("Query error: %s", err)
		return VCStatusErrorResponse
	}
	if n >= len(data) {
		io.Copy(ioutil.Discard, resp.Body)
	}
	resp.Body.Close()

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "text/xml") {
		op.Errorf("Unexpected content type %s, should be text/xml", contentType)
		op.Errorf("Response from the server: %s", string(data))
		return VCStatusNotXML
	}
	// we just want to make sure that response contains something familiar that we could
	// user as ESXi marker.
	if !bytes.Contains(data, []byte("urn:vim25Service")) {
		op.Errorf("Server response doesn't contain 'urn:vim25Service': %s", string(data))
		return VCStatusIncorrectResponse
	}
	return VCStatusOK
}
