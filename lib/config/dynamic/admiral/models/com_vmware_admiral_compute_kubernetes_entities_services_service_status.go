package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus com vmware admiral compute kubernetes entities services service status
// swagger:model com:vmware:admiral:compute:kubernetes:entities:services:ServiceStatus
type ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus struct {

	// load balancer
	LoadBalancer *ComVmwareAdmiralComputeKubernetesEntitiesServicesLoadBalancerStatus `json:"loadBalancer,omitempty"`
}

// Validate validates this com vmware admiral compute kubernetes entities services service status
func (m *ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLoadBalancer(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus) validateLoadBalancer(formats strfmt.Registry) error {

	if swag.IsZero(m.LoadBalancer) { // not required
		return nil
	}

	if m.LoadBalancer != nil {

		if err := m.LoadBalancer.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("loadBalancer")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralComputeKubernetesEntitiesServicesServiceStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
