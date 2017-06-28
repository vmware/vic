package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec com vmware admiral compute kubernetes entities pods pod template spec
// swagger:model com:vmware:admiral:compute:kubernetes:entities:pods:PodTemplateSpec
type ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec struct {

	// metadata
	Metadata *ComVmwareAdmiralComputeKubernetesEntitiesCommonObjectMeta `json:"metadata,omitempty"`

	// spec
	Spec *ComVmwareAdmiralComputeKubernetesEntitiesPodsPodSpec `json:"spec,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities pods pod template spec
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetadata(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateSpec(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec) validateMetadata(formats strfmt.Registry) error {

	if swag.IsZero(m.Metadata) { // not required
		return nil
	}

	if m.Metadata != nil {

		if err := m.Metadata.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metadata")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec) validateSpec(formats strfmt.Registry) error {

	if swag.IsZero(m.Spec) { // not required
		return nil
	}

	if m.Spec != nil {

		if err := m.Spec.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("spec")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesPodsPodTemplateSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
