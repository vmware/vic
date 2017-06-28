package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource com vmware admiral compute kubernetes entities volumes flex volume source
// swagger:model com:vmware:admiral:compute:kubernetes:entities:volumes:FlexVolumeSource
type ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource struct {

	// driver
	Driver string `json:"driver,omitempty"`

	// fs type
	FsType string `json:"fsType,omitempty"`

	// options
	Options map[string]string `json:"options,omitempty"`

	// read only
	ReadOnly bool `json:"readOnly,omitempty"`

	// secret ref
	SecretRef *ComVmwareAdmiralComputeKubernetesEntitiesCommonLocalObjectReference `json:"secretRef,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities volumes flex volume source
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSecretRef(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource) validateSecretRef(formats strfmt.Registry) error {

	if swag.IsZero(m.SecretRef) { // not required
		return nil
	}

	if m.SecretRef != nil {

		if err := m.SecretRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("secretRef")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesVolumesFlexVolumeSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
