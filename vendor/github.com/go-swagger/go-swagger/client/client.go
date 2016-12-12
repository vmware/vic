// Copyright 2015 go-swagger maintainers
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

package client

import "golang.org/x/net/context"

// Operation represents the context for a swagger operation to be submitted to the transport
type Operation struct {
	ID                 string
	Method             string
	PathPattern        string
	ProducesMediaTypes []string
	ConsumesMediaTypes []string
	Schemes            []string
	AuthInfo           AuthInfoWriter
	Params             RequestWriter
	Reader             ResponseReader
	Context            context.Context
}

// A Transport implementor knows how to submit Request objects to some destination
type Transport interface {
	//Submit(string, RequestWriter, ResponseReader, AuthInfoWriter) (interface{}, error)
	Submit(*Operation) (interface{}, error)
}
