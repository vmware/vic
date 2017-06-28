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

// ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState com vmware admiral compute container network container network service container network state
// swagger:model com:vmware:admiral:compute:container:network:ContainerNetworkService:ContainerNetworkState
type ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState struct {

	// Defines which adapter will serve the provision request
	AdapterManagementReference strfmt.URI `json:"adapterManagementReference,omitempty"`

	// Links to CompositeComponents when a network is part of App/Composition request.
	CompositeComponentLinks []string `json:"compositeComponentLinks"`

	// Network connected time in milliseconds
	Connected int64 `json:"connected,omitempty"`

	// Runtime property that will be populated during network inspections. Contains the number of containers that are connected to this container network.
	ConnectedContainersCount int64 `json:"connectedContainersCount,omitempty"`

	// creation time micros
	CreationTimeMicros int64 `json:"creationTimeMicros,omitempty"`

	// custom properties
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// desc
	Desc string `json:"desc,omitempty"`

	// Defines the description of the network.
	DescriptionLink string `json:"descriptionLink,omitempty"`

	// The name of the driver for this network. Can be bridge, host, overlay, none.
	Driver string `json:"driver,omitempty"`

	// If set to true, specifies that this network exists independently of any application.
	External bool `json:"external,omitempty"`

	// group links
	GroupLinks []string `json:"groupLinks"`

	// id
	ID string `json:"id,omitempty"`

	// An IPAM configuration for a given network.
	IPAM *ComVmwareAdmiralComputeContainerNetworkIPAM `json:"ipam,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// A map of field-value pairs for a given network. These are used to specify network options that are used by the network drivers.
	Options map[string]string `json:"options,omitempty"`

	// Reference to the host that this network was created on.
	OriginatingHostLink string `json:"originatingHostLink,omitempty"`

	// Container host links
	ParentLinks []string `json:"parentLinks"`

	// Network state indicating runtime state of a network instance.
	PowerState string `json:"powerState,omitempty"`

	// region Id
	RegionID string `json:"regionId,omitempty"`

	// tag links
	TagLinks []string `json:"tagLinks"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware admiral compute container network container network service container network state
func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCompositeComponentLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateGroupLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateIPAM(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateParentLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validatePowerState(formats); err != nil {
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

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validateCompositeComponentLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.CompositeComponentLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validateGroupLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.GroupLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validateIPAM(formats strfmt.Registry) error {

	if swag.IsZero(m.IPAM) { // not required
		return nil
	}

	if m.IPAM != nil {

		if err := m.IPAM.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("ipam")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validateParentLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.ParentLinks) { // not required
		return nil
	}

	return nil
}

var comVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStateTypePowerStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["UNKNOWN","PROVISIONING","CONNECTED","RETIRED","ERROR"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		comVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStateTypePowerStatePropEnum = append(comVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStateTypePowerStatePropEnum, v)
	}
}

const (
	// ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateUNKNOWN captures enum value "UNKNOWN"
	ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateUNKNOWN string = "UNKNOWN"
	// ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStatePROVISIONING captures enum value "PROVISIONING"
	ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStatePROVISIONING string = "PROVISIONING"
	// ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateCONNECTED captures enum value "CONNECTED"
	ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateCONNECTED string = "CONNECTED"
	// ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateRETIRED captures enum value "RETIRED"
	ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateRETIRED string = "RETIRED"
	// ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateERROR captures enum value "ERROR"
	ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStatePowerStateERROR string = "ERROR"
)

// prop value enum
func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validatePowerStateEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, comVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkStateTypePowerStatePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validatePowerState(formats strfmt.Registry) error {

	if swag.IsZero(m.PowerState) { // not required
		return nil
	}

	// value enum
	if err := m.validatePowerStateEnum("powerState", "body", m.PowerState); err != nil {
		return err
	}

	return nil
}

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validateTagLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TagLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeContainerNetworkContainerNetworkServiceContainerNetworkState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
