package main

import (
	_ "github.com/docker/swarmkit/protobuf/plugin/authenticatedwrapper"
	_ "github.com/docker/swarmkit/protobuf/plugin/deepcopy"
	_ "github.com/docker/swarmkit/protobuf/plugin/raftproxy"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()
	files = vanity.FilterFiles(files, vanity.NotInPackageGoogleProtobuf)

	for _, opt := range []func(*descriptor.FileDescriptorProto){
		vanity.TurnOnGoStringAll,
		vanity.TurnOffGoGettersAll,
		vanity.TurnOffGoStringerAll,
		vanity.TurnOnMarshalerAll,
		vanity.TurnOnStringerAll,
		vanity.TurnOnUnmarshalerAll,
		vanity.TurnOnSizerAll,
	} {
		vanity.ForEachFile(files, opt)
	}

	resp := command.Generate(req)
	command.Write(resp)
}
