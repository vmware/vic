package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState com vmware admiral compute container container host data collection service container host data collection state
// swagger:model com:vmware:admiral:compute:container:ContainerHostDataCollectionService:ContainerHostDataCollectionState
type ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState struct {

	// List of container host links to be updated as part of Patch triggered data collection.
	ComputeContainerHostLinks []string `json:"computeContainerHostLinks"`

	// create or update host
	CreateOrUpdateHost bool `json:"createOrUpdateHost,omitempty"`

	// Indicator of the last run of data-collection.
	LastRunTimeMicros int64 `json:"lastRunTimeMicros,omitempty"`

	// Flag indicating if this is data-collection after container remove.
	Remove bool `json:"remove,omitempty"`

	// Count of how many times the last run data-collection has been run within very small time period.
	SkipRunCount int64 `json:"skipRunCount,omitempty"`
}

// Validate validates this com vmware admiral compute container container host data collection service container host data collection state
func (m *ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComputeContainerHostLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState) validateComputeContainerHostLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.ComputeContainerHostLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeContainerContainerHostDataCollectionServiceContainerHostDataCollectionState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
