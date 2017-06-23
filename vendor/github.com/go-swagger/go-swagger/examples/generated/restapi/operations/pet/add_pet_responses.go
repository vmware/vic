package pet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// AddPetMethodNotAllowedCode is the HTTP code returned for type AddPetMethodNotAllowed
const AddPetMethodNotAllowedCode int = 405

/*AddPetMethodNotAllowed Invalid input

swagger:response addPetMethodNotAllowed
*/
type AddPetMethodNotAllowed struct {
}

// NewAddPetMethodNotAllowed creates AddPetMethodNotAllowed with default headers values
func NewAddPetMethodNotAllowed() *AddPetMethodNotAllowed {
	return &AddPetMethodNotAllowed{}
}

// WriteResponse to the client
func (o *AddPetMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
}
