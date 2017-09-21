package vchlog

import (
	"sync"
	"bytes"
)

type BufferedPipe struct {
	buffer *bytes.Buffer
	c *sync.Cond
	closed bool
}

func NewBufferedPipe() *BufferedPipe {
	var m sync.Mutex
	c := sync.NewCond(&m)
	return &BufferedPipe{
		buffer: bytes.NewBuffer(nil),
		c: c,
		closed: false,
	}
}

func (bp *BufferedPipe) Read(data []byte) (n int, err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()

	for bp.buffer.Len() == 0 && !bp.closed {
		bp.c.Wait()
	}

	return bp.buffer.Read(data)
}

func (bp *BufferedPipe) Write(data []byte) (n int, err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Signal()
	return bp.buffer.Write(data)
}

func (bp *BufferedPipe) Close(err error) {
	bp.c.L.Lock()
	defer bp.c.L.Unlock()
	defer bp.c.Signal()
	bp.closed = true
}