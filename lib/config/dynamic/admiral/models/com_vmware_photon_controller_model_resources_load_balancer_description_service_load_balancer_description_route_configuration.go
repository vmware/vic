package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration com vmware photon controller model resources load balancer description service load balancer description route configuration
// swagger:model com:vmware:photon:controller:model:resources:LoadBalancerDescriptionService:LoadBalancerDescription:RouteConfiguration
type ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration struct {

	// health check configuration
	HealthCheckConfiguration *ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionHealthCheckConfiguration `json:"healthCheckConfiguration,omitempty"`

	// instance port
	InstancePort string `json:"instancePort,omitempty"`

	// instance protocol
	InstanceProtocol string `json:"instanceProtocol,omitempty"`

	// port
	Port string `json:"port,omitempty"`

	// protocol
	Protocol string `json:"protocol,omitempty"`
}

// Validate validates this com vmware photon controller model resources load balancer description service load balancer description route configuration
func (m *ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHealthCheckConfiguration(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration) validateHealthCheckConfiguration(formats strfmt.Registry) error {

	if swag.IsZero(m.HealthCheckConfiguration) { // not required
		return nil
	}

	if m.HealthCheckConfiguration != nil {

		if err := m.HealthCheckConfiguration.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("healthCheckConfiguration")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration) UnmarshalBinary(b []byte) error {
	var res ComVmwarePhotonControllerModelResourcesLoadBalancerDescriptionServiceLoadBalancerDescriptionRouteConfiguration
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
