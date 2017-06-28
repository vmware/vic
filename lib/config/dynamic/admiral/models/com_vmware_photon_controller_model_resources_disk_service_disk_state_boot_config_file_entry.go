package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwarePhotonControllerModelResourcesDiskServiceDiskStateBootConfigFileEntry com vmware photon controller model resources disk service disk state boot config file entry
// swagger:model com:vmware:photon:controller:model:resources:DiskService:DiskState:BootConfig:FileEntry
type ComVmwarePhotonControllerModelResourcesDiskServiceDiskStateBootConfigFileEntry struct {

	// contents
	Contents string `json:"contents,omitempty"`

	// contents reference
	ContentsReference strfmt.URI `json:"contentsReference,omitempty"`

	// path
	Path string `json:"path,omitempty"`
}

// Validate validates this com vmware photon controller model resources disk service disk state boot config file entry
func (m *ComVmwarePhotonControllerModelResourcesDiskServiceDiskStateBootConfigFileEntry) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesDiskServiceDiskStateBootConfigFileEntry) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesDiskServiceDiskStateBootConfigFileEntry) UnmarshalBinary(b []byte) error {
	var res ComVmwarePhotonControllerModelResourcesDiskServiceDiskStateBootConfigFileEntry
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
