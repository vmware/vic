package main

import (
	"log"

	spec "github.com/go-swagger/go-swagger/spec"
	flags "github.com/jessevdk/go-flags"

	"github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-1/restapi"
	"github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-1/restapi/operations"
)

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	swaggerSpec, err := spec.New(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTodoListAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer api.ServerShutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = swaggerSpec.Spec().Info.Title
	parser.LongDescription = swaggerSpec.Spec().Info.Description

	for _, optsGroup := range api.CommandLineOptionsGroups {
		parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
	}

	if _, err := parser.Parse(); err != nil {
		log.Fatalln(err)
	}

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
