package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesVolumesCinderVolumeSource com vmware admiral compute kubernetes entities volumes cinder volume source
// swagger:model com:vmware:admiral:compute:kubernetes:entities:volumes:CinderVolumeSource
type ComVmwareAdmiralComputeKubernetesEntitiesVolumesCinderVolumeSource struct {

	// fs type
	FsType string `json:"fsType,omitempty"`

	// read only
	ReadOnly bool `json:"readOnly,omitempty"`

	// volume ID
	VolumeID string `json:"volumeID,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities volumes cinder volume source
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesCinderVolumeSource) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesCinderVolumeSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesCinderVolumeSource) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesVolumesCinderVolumeSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
