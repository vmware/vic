package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState com vmware admiral request composition composition sub task service composition sub task state
// swagger:model com:vmware:admiral:request:composition:CompositionSubTaskService:CompositionSubTaskState
type ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState struct {

	// allocation request
	AllocationRequest bool `json:"allocationRequest,omitempty"`

	// composite description link
	CompositeDescriptionLink string `json:"compositeDescriptionLink,omitempty"`

	// current depends on link
	CurrentDependsOnLink string `json:"currentDependsOnLink,omitempty"`

	// Custom properties.
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// dependent links
	DependentLinks []string `json:"dependentLinks"`

	// depends on links
	DependsOnLinks []string `json:"dependsOnLinks"`

	// error count
	ErrorCount int64 `json:"errorCount,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// operation
	Operation string `json:"operation,omitempty"`

	// post allocation
	PostAllocation bool `json:"postAllocation,omitempty"`

	// request Id
	RequestID string `json:"requestId,omitempty"`

	// link to a service that will receive updates when the task changes state.
	RequestTrackerLink string `json:"requestTrackerLink,omitempty"`

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

// Validate validates this com vmware admiral request composition composition sub task service composition sub task state
func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDependentLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDependsOnLinks(formats); err != nil {
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

func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) validateDependentLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.DependentLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) validateDependsOnLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.DependsOnLinks) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) validateServiceTaskCallback(formats strfmt.Registry) error {

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

func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) validateTaskInfo(formats strfmt.Registry) error {

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

func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralRequestCompositionCompositionSubTaskServiceCompositionSubTaskState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
