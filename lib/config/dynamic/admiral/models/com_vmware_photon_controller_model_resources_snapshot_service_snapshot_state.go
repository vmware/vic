package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState com vmware photon controller model resources snapshot service snapshot state
// swagger:model com:vmware:photon:controller:model:resources:SnapshotService:SnapshotState
type ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState struct {

	// compute link
	ComputeLink string `json:"computeLink,omitempty"`

	// creation time micros
	CreationTimeMicros int64 `json:"creationTimeMicros,omitempty"`

	// custom properties
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// desc
	Desc string `json:"desc,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// group links
	GroupLinks []string `json:"groupLinks"`

	// id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// region Id
	RegionID string `json:"regionId,omitempty"`

	// tag links
	TagLinks []string `json:"tagLinks"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware photon controller model resources snapshot service snapshot state
func (m *ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGroupLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTagLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTenantLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState) validateGroupLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.GroupLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState) validateTagLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TagLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState) UnmarshalBinary(b []byte) error {
	var res ComVmwarePhotonControllerModelResourcesSnapshotServiceSnapshotState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
