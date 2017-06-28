package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesPodsTCPSocketAction com vmware admiral compute kubernetes entities pods TCP socket action
// swagger:model com:vmware:admiral:compute:kubernetes:entities:pods:TCPSocketAction
type ComVmwareAdmiralComputeKubernetesEntitiesPodsTCPSocketAction struct {

	// port
	Port string `json:"port,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities pods TCP socket action
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsTCPSocketAction) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsTCPSocketAction) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsTCPSocketAction) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesPodsTCPSocketAction
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
