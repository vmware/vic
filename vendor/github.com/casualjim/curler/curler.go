package curler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// New creates a new curler
func New(orig http.Handler, out io.Writer) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		parts := []string{"curl", "-i", "-X", r.Method}

		for k, v := range r.Header {
			ck := http.CanonicalHeaderKey(k)
			if ck != "Host" && ck != "User-Agent" {
				parts = append(parts, "-H", fmt.Sprintf("'%s: %s'", k, strings.Join(v, ",")))
			}
		}

		if r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" {
			b, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Fprintln(out, "curler:", err)
			}
			if err == nil && len(b) > 0 {
				r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
				parts = append(parts, "-d '", string(b)+"'")
			}
		}

		scheme := "http"
		if r.TLS != nil {
			scheme += "s"
		}

		parts = append(parts, scheme+"://"+r.Host+r.URL.String())
		fmt.Fprintln(out, strings.Join(parts, " "))
		orig.ServeHTTP(rw, r)
	})
}
