package diag

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/trace"
)

func TestPing(t *testing.T) {
	assert.Equal(t, PingStatusNoPingOutput, runPing(nil, nil))
	assert.Equal(t, PingStatusPingNotExists, runPing(nil, errors.New("executable file not found")))
	assert.Equal(t, PingStatusResolutionFailed, runPing([]byte("unknown host test"), nil))
	assert.Equal(t, PingStatusUnknownError, runPing([]byte("strageoutput"), nil))

	pingOutput1 := `PING ukr.net (212.42.76.253) 56(84) bytes of data.
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=1 ttl=47 time=168 ms
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=2 ttl=47 time=173 ms
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=3 ttl=47 time=172 ms
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=4 ttl=47 time=164 ms

--- ukr.net ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 164.722/169.911/173.932/3.512 ms
	`

	assert.Equal(t, PingStatusOk, runPing([]byte(pingOutput1), nil))

	pingOutput2 := `PING 129.1.0.8 (129.1.0.8) 56(84) bytes of data.

--- 129.1.0.8 ping statistics ---
4 packets transmitted, 0 received, 100% packet loss, time 3024ms
`
	assert.Equal(t, PingStatusOkNotPingable, runPing([]byte(pingOutput2), nil))

	pingOutput3 := `PING ukr.net (212.42.76.253) 56(84) bytes of data.
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=1 ttl=47 time=168 ms
64 bytes from srv253.fwdcdn.com (212.42.76.253): timeout
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=3 ttl=47 time=172 ms
64 bytes from srv253.fwdcdn.com (212.42.76.253): icmp_seq=4 ttl=47 time=164 ms

--- ukr.net ping statistics ---
4 packets transmitted, 3 received, 25% packet loss, time 3005ms
rtt min/avg/max/mdev = 164.722/169.911/173.932/3.512 ms
	`

	assert.Equal(t, PingStatusOkPacketLosses, runPing([]byte(pingOutput3), nil))
}

func TestCheckAPIAvailability(t *testing.T) {
	assert.Equal(t, VCStatusErrorQuery, CheckAPIAvailability("http://127.0.0.1:65535"))
	assert.Equal(t, VCStatusErrorQuery, CheckAPIAvailability("http://127.0.0.1:65536"))
}

func TestCheckAPIAvailabilityQueryWithGetterError(t *testing.T) {
	op := trace.NewOperation(context.Background(), "test")
	f := func(s string) (*http.Response, error) { return nil, errors.New("wrong query") }
	code := queryAPI(op, f, "testurl")
	assert.Equal(t, VCStatusErrorQuery, code)
}

type readerWithError struct {
	err  error
	data *bytes.Reader
}

func (r *readerWithError) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.data.Read(b)
}

func (r *readerWithError) Close() error {
	return r.err
}

func TestCheckAPIAvailabilityQueryReadError(t *testing.T) {
	op := trace.NewOperation(context.Background(), "test")
	f := func(s string) (*http.Response, error) {
		hr := &http.Response{
			Body: &readerWithError{
				err: errors.New("read error happened"),
			},
		}
		return hr, nil
	}
	code := queryAPI(op, f, "testurl")
	assert.Equal(t, VCStatusErrorResponse, code)
}

func TestCheckAPIAvailabilityQueryIncorrectDataType(t *testing.T) {
	op := trace.NewOperation(context.Background(), "test")
	f := func(s string) (*http.Response, error) {
		hr := &http.Response{
			Body: &readerWithError{
				data: bytes.NewReader([]byte("some data")),
			},
		}
		return hr, nil
	}
	code := queryAPI(op, f, "testurl")
	assert.Equal(t, VCStatusNotXML, code)
}

func TestCheckAPIAvailabilityQueryIncorrectData(t *testing.T) {
	op := trace.NewOperation(context.Background(), "test")
	f := func(s string) (*http.Response, error) {
		hr := &http.Response{
			Body: &readerWithError{
				data: bytes.NewReader([]byte("some data")),
			},
			Header: http.Header{"Content-Type": []string{"text/xml"}},
		}
		return hr, nil
	}
	code := queryAPI(op, f, "testurl")
	assert.Equal(t, VCStatusIncorrectResponse, code)
}

func TestCheckAPIAvailabilityQueryCorrectData(t *testing.T) {
	op := trace.NewOperation(context.Background(), "test")
	f := func(s string) (*http.Response, error) {
		hr := &http.Response{
			Body: &readerWithError{
				data: bytes.NewReader([]byte("some urn:vim25Service data")),
			},
			Header: http.Header{"Content-Type": []string{"text/xml"}},
		}
		return hr, nil
	}
	code := queryAPI(op, f, "testurl")
	assert.Equal(t, VCStatusOK, code)
}
