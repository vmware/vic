package xdr

import (
	"io"

	xdr "github.com/davecgh/go-xdr/xdr2"
)

func Write(w io.Writer, val interface{}) error {
	_, err := xdr.Marshal(w, val)
	return err
}
