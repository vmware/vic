package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState com vmware admiral request request broker service request broker state
// swagger:model com:vmware:admiral:request:RequestBrokerService:RequestBrokerState
type ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState struct {

	// Custom properties.
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// operation
	Operation string `json:"operation,omitempty"`

	// link to a service that will receive updates when the task changes state.
	RequestTrackerLink string `json:"requestTrackerLink,omitempty"`

	// resource count
	ResourceCount int64 `json:"resourceCount,omitempty"`

	// resource description link
	ResourceDescriptionLink string `json:"resourceDescriptionLink,omitempty"`

	// resource type
	ResourceType string `json:"resourceType,omitempty"`

	// Callback link and response from the service initiated this task.
	ServiceTaskCallback *ComVmwareAdmiralServiceCommonServiceTaskCallback `json:"serviceTaskCallback,omitempty"`

	//  Describes a service task state.
	TaskInfo *ComVmwareXenonCommonTaskState `json:"taskInfo,omitempty"`

	//  Describes a service task sub stage.
	TaskSubStage string `json:"taskSubStage,omitempty"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware admiral request request broker service request broker state
func (m *ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateServiceTaskCallback(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTaskInfo(formats); err != nil {
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

func (m *ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState) validateServiceTaskCallback(formats strfmt.Registry) error {

	if swag.IsZero(m.ServiceTaskCallback) { // not required
		return nil
	}

	if m.ServiceTaskCallback != nil {

		if err := m.ServiceTaskCallback.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("serviceTaskCallback")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState) validateTaskInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.TaskInfo) { // not required
		return nil
	}

	if m.TaskInfo != nil {

		if err := m.TaskInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("taskInfo")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralRequestRequestBrokerServiceRequestBrokerState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
