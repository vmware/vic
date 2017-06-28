package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesVolumesAzureDiskVolumeSource com vmware admiral compute kubernetes entities volumes azure disk volume source
// swagger:model com:vmware:admiral:compute:kubernetes:entities:volumes:AzureDiskVolumeSource
type ComVmwareAdmiralComputeKubernetesEntitiesVolumesAzureDiskVolumeSource struct {

	// caching mode
	CachingMode JavaLangObject `json:"cachingMode,omitempty"`

	// disk name
	DiskName string `json:"diskName,omitempty"`

	// disk URI
	DiskURI string `json:"diskURI,omitempty"`

	// fs type
	FsType string `json:"fsType,omitempty"`

	// read only
	ReadOnly bool `json:"readOnly,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities volumes azure disk volume source
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesAzureDiskVolumeSource) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesAzureDiskVolumeSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesAzureDiskVolumeSource) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesVolumesAzureDiskVolumeSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
