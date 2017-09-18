package vchlog

import (
	"github.com/vmware/govmomi/object"
	"context"
	"path"
)


type VCHCreatedSignal struct {
	Datastore *object.Datastore
	LogFileName string
	Contex context.Context
	VMPathName string
}

var Pipe *BufferedPipe
var SignalChan chan VCHCreatedSignal


func Init() {
	Pipe = NewBufferedPipe()
	SignalChan = make(chan VCHCreatedSignal)
}

func Run() {
	sig := <-SignalChan
	sig.Datastore.Upload(sig.Contex, Pipe, path.Join(sig.VMPathName, sig.LogFileName), nil)
}

func GetPipe() *BufferedPipe {
	return Pipe
}

func Signal(sig VCHCreatedSignal) {
	SignalChan <- sig
}





