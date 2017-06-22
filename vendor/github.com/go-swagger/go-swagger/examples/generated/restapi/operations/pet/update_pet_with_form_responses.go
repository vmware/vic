package pet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// UpdatePetWithFormMethodNotAllowedCode is the HTTP code returned for type UpdatePetWithFormMethodNotAllowed
const UpdatePetWithFormMethodNotAllowedCode int = 405

/*UpdatePetWithFormMethodNotAllowed Invalid input

swagger:response updatePetWithFormMethodNotAllowed
*/
type UpdatePetWithFormMethodNotAllowed struct {
}

// NewUpdatePetWithFormMethodNotAllowed creates UpdatePetWithFormMethodNotAllowed with default headers values
func NewUpdatePetWithFormMethodNotAllowed() *UpdatePetWithFormMethodNotAllowed {
	return &UpdatePetWithFormMethodNotAllowed{}
}

// WriteResponse to the client
func (o *UpdatePetWithFormMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
}
