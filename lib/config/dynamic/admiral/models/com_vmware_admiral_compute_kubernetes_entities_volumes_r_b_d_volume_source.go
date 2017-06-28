package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource com vmware admiral compute kubernetes entities volumes r b d volume source
// swagger:model com:vmware:admiral:compute:kubernetes:entities:volumes:RBDVolumeSource
type ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource struct {

	// fs type
	FsType string `json:"fsType,omitempty"`

	// image
	Image string `json:"image,omitempty"`

	// keyring
	Keyring string `json:"keyring,omitempty"`

	// monitors
	Monitors []string `json:"monitors"`

	// pool
	Pool string `json:"pool,omitempty"`

	// read only
	ReadOnly bool `json:"readOnly,omitempty"`

	// source ref
	SourceRef *ComVmwareAdmiralComputeKubernetesEntitiesCommonLocalObjectReference `json:"sourceRef,omitempty"`

	// user
	User string `json:"user,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities volumes r b d volume source
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMonitors(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateSourceRef(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource) validateMonitors(formats strfmt.Registry) error {

	if swag.IsZero(m.Monitors) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource) validateSourceRef(formats strfmt.Registry) error {

	if swag.IsZero(m.SourceRef) { // not required
		return nil
	}

	if m.SourceRef != nil {

		if err := m.SourceRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sourceRef")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesVolumesRBDVolumeSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
