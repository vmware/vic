package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState com vmware admiral compute elastic placement zone configuration service elastic placement zone configuration state
// swagger:model com:vmware:admiral:compute:ElasticPlacementZoneConfigurationService:ElasticPlacementZoneConfigurationState
type ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState struct {

	// epz state
	EpzState *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState `json:"epzState,omitempty"`

	// resource pool state
	ResourcePoolState *ComVmwarePhotonControllerModelResourcesResourcePoolServiceResourcePoolState `json:"resourcePoolState,omitempty"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware admiral compute elastic placement zone configuration service elastic placement zone configuration state
func (m *ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEpzState(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateResourcePoolState(formats); err != nil {
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

func (m *ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState) validateEpzState(formats strfmt.Registry) error {

	if swag.IsZero(m.EpzState) { // not required
		return nil
	}

	if m.EpzState != nil {

		if err := m.EpzState.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("epzState")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState) validateResourcePoolState(formats strfmt.Registry) error {

	if swag.IsZero(m.ResourcePoolState) { // not required
		return nil
	}

	if m.ResourcePoolState != nil {

		if err := m.ResourcePoolState.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("resourcePoolState")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeElasticPlacementZoneConfigurationServiceElasticPlacementZoneConfigurationState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
