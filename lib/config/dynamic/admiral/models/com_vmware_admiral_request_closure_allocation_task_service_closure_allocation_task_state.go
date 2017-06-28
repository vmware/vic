package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState com vmware admiral request closure allocation task service closure allocation task state
// swagger:model com:vmware:admiral:request:ClosureAllocationTaskService:ClosureAllocationTaskState
type ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState struct {

	// The description that defines the closure.
	ClosureLink string `json:"closureLink,omitempty"`

	// Custom properties.
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// link to a service that will receive updates when the task changes state.
	RequestTrackerLink string `json:"requestTrackerLink,omitempty"`

	// The description that defines the closure description.
	ResourceDescriptionLink string `json:"resourceDescriptionLink,omitempty"`

	// Set by a Task with the links of the provisioned resources.
	ResourceLinks []string `json:"resourceLinks"`

	// Callback link and response from the service initiated this task.
	ServiceTaskCallback *ComVmwareAdmiralServiceCommonServiceTaskCallback `json:"serviceTaskCallback,omitempty"`

	//  Describes a service task state.
	TaskInfo *ComVmwareXenonCommonTaskState `json:"taskInfo,omitempty"`

	//  Describes a service task sub stage.
	TaskSubStage string `json:"taskSubStage,omitempty"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`
}

// Validate validates this com vmware admiral request closure allocation task service closure allocation task state
func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateResourceLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

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

func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) validateResourceLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.ResourceLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) validateServiceTaskCallback(formats strfmt.Registry) error {

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

func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) validateTaskInfo(formats strfmt.Registry) error {

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

func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralRequestClosureAllocationTaskServiceClosureAllocationTaskState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
