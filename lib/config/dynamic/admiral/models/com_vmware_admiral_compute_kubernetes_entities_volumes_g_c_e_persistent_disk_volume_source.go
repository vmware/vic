package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesVolumesGCEPersistentDiskVolumeSource com vmware admiral compute kubernetes entities volumes g c e persistent disk volume source
// swagger:model com:vmware:admiral:compute:kubernetes:entities:volumes:GCEPersistentDiskVolumeSource
type ComVmwareAdmiralComputeKubernetesEntitiesVolumesGCEPersistentDiskVolumeSource struct {

	// fs type
	FsType string `json:"fsType,omitempty"`

	// partition
	Partition int64 `json:"partition,omitempty"`

	// pd name
	PdName string `json:"pdName,omitempty"`

	// read only
	ReadOnly bool `json:"readOnly,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities volumes g c e persistent disk volume source
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesGCEPersistentDiskVolumeSource) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesGCEPersistentDiskVolumeSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesVolumesGCEPersistentDiskVolumeSource) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesVolumesGCEPersistentDiskVolumeSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
