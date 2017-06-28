package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesPodsHTTPGetAction com vmware admiral compute kubernetes entities pods HTTP get action
// swagger:model com:vmware:admiral:compute:kubernetes:entities:pods:HTTPGetAction
type ComVmwareAdmiralComputeKubernetesEntitiesPodsHTTPGetAction struct {

	// host
	Host string `json:"host,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// port
	Port string `json:"port,omitempty"`

	// scheme
	Scheme string `json:"scheme,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities pods HTTP get action
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsHTTPGetAction) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsHTTPGetAction) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsHTTPGetAction) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesPodsHTTPGetAction
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
