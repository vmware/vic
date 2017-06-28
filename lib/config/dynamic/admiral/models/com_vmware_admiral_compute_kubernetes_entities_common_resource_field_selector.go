package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesCommonResourceFieldSelector com vmware admiral compute kubernetes entities common resource field selector
// swagger:model com:vmware:admiral:compute:kubernetes:entities:common:ResourceFieldSelector
type ComVmwareAdmiralComputeKubernetesEntitiesCommonResourceFieldSelector struct {

	// container name
	ContainerName string `json:"containerName,omitempty"`

	// divisor
	Divisor string `json:"divisor,omitempty"`

	// resource
	Resource string `json:"resource,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities common resource field selector
func (m *ComVmwareAdmiralComputeKubernetesEntitiesCommonResourceFieldSelector) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesCommonResourceFieldSelector) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesCommonResourceFieldSelector) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesCommonResourceFieldSelector
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
