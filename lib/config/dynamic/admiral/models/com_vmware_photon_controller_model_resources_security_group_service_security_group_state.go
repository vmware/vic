package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState com vmware photon controller model resources security group service security group state
// swagger:model com:vmware:photon:controller:model:resources:SecurityGroupService:SecurityGroupState
type ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState struct {

	// auth credentials link
	AuthCredentialsLink string `json:"authCredentialsLink,omitempty"`

	// creation time micros
	CreationTimeMicros int64 `json:"creationTimeMicros,omitempty"`

	// custom properties
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// desc
	Desc string `json:"desc,omitempty"`

	// egress
	Egress []*ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupStateRule `json:"egress"`

	// endpoint link
	EndpointLink string `json:"endpointLink,omitempty"`

	// group links
	GroupLinks []string `json:"groupLinks"`

	// id
	ID string `json:"id,omitempty"`

	// ingress
	Ingress []*ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupStateRule `json:"ingress"`

	// instance adapter reference
	InstanceAdapterReference strfmt.URI `json:"instanceAdapterReference,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// region Id
	RegionID string `json:"regionId,omitempty"`

	// resource pool link
	ResourcePoolLink string `json:"resourcePoolLink,omitempty"`

	// tag links
	TagLinks []string `json:"tagLinks"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware photon controller model resources security group service security group state
func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEgress(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateGroupLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateIngress(formats); err != nil {
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

func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) validateEgress(formats strfmt.Registry) error {

	if swag.IsZero(m.Egress) { // not required
		return nil
	}

	for i := 0; i < len(m.Egress); i++ {

		if swag.IsZero(m.Egress[i]) { // not required
			continue
		}

		if m.Egress[i] != nil {

			if err := m.Egress[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("egress" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) validateGroupLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.GroupLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) validateIngress(formats strfmt.Registry) error {

	if swag.IsZero(m.Ingress) { // not required
		return nil
	}

	for i := 0; i < len(m.Ingress); i++ {

		if swag.IsZero(m.Ingress[i]) { // not required
			continue
		}

		if m.Ingress[i] != nil {

			if err := m.Ingress[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ingress" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) validateTagLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TagLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState) UnmarshalBinary(b []byte) error {
	var res ComVmwarePhotonControllerModelResourcesSecurityGroupServiceSecurityGroupState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
