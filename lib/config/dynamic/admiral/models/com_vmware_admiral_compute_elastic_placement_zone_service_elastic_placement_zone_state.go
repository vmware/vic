package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState com vmware admiral compute elastic placement zone service elastic placement zone state
// swagger:model com:vmware:admiral:compute:ElasticPlacementZoneService:ElasticPlacementZoneState
type ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState struct {

	// Advanced placement policy.
	PlacementPolicy string `json:"placementPolicy,omitempty"`

	// Link to the elastic resource pool
	ResourcePoolLink string `json:"resourcePoolLink,omitempty"`

	// Links to tags that must be set on the computes in order to add them to the elastic resource pool
	TagLinksToMatch []string `json:"tagLinksToMatch"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware admiral compute elastic placement zone service elastic placement zone state
func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePlacementPolicy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTagLinksToMatch(formats); err != nil {
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

var comVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStateTypePlacementPolicyPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["DEFAULT","SPREAD","BINPACK"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		comVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStateTypePlacementPolicyPropEnum = append(comVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStateTypePlacementPolicyPropEnum, v)
	}
}

const (
	// ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStatePlacementPolicyDEFAULT captures enum value "DEFAULT"
	ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStatePlacementPolicyDEFAULT string = "DEFAULT"
	// ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStatePlacementPolicySPREAD captures enum value "SPREAD"
	ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStatePlacementPolicySPREAD string = "SPREAD"
	// ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStatePlacementPolicyBINPACK captures enum value "BINPACK"
	ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStatePlacementPolicyBINPACK string = "BINPACK"
)

// prop value enum
func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) validatePlacementPolicyEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, comVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneStateTypePlacementPolicyPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) validatePlacementPolicy(formats strfmt.Registry) error {

	if swag.IsZero(m.PlacementPolicy) { // not required
		return nil
	}

	// value enum
	if err := m.validatePlacementPolicyEnum("placementPolicy", "body", m.PlacementPolicy); err != nil {
		return err
	}

	return nil
}

func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) validateTagLinksToMatch(formats strfmt.Registry) error {

	if swag.IsZero(m.TagLinksToMatch) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeElasticPlacementZoneServiceElasticPlacementZoneState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
