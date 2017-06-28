package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ComVmwareXenonServicesCommonMigrationTaskServiceState com vmware xenon services common migration task service state
// swagger:model com:vmware:xenon:services:common:MigrationTaskService:State
type ComVmwareXenonServicesCommonMigrationTaskServiceState struct {

	// continuous migration
	ContinuousMigration bool `json:"continuousMigration,omitempty"`

	// destination factory link
	DestinationFactoryLink string `json:"destinationFactoryLink,omitempty"`

	// destination node group reference
	DestinationNodeGroupReference strfmt.URI `json:"destinationNodeGroupReference,omitempty"`

	// latest source update time micros
	LatestSourceUpdateTimeMicros int64 `json:"latestSourceUpdateTimeMicros,omitempty"`

	// maintenance interval micros
	MaintenanceIntervalMicros int64 `json:"maintenanceIntervalMicros,omitempty"`

	// maximum convergence checks
	MaximumConvergenceChecks int64 `json:"maximumConvergenceChecks,omitempty"`

	// migration options
	MigrationOptions []string `json:"migrationOptions"`

	// query spec
	QuerySpec *ComVmwareXenonServicesCommonQueryTaskQuerySpecification `json:"querySpec,omitempty"`

	// source factory link
	SourceFactoryLink string `json:"sourceFactoryLink,omitempty"`

	// source node group reference
	SourceNodeGroupReference strfmt.URI `json:"sourceNodeGroupReference,omitempty"`

	// task info
	TaskInfo *ComVmwareXenonCommonTaskState `json:"taskInfo,omitempty"`

	// transformation service link
	TransformationServiceLink string `json:"transformationServiceLink,omitempty"`
}

// Validate validates this com vmware xenon services common migration task service state
func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMigrationOptions(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateQuerySpec(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTaskInfo(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var comVmwareXenonServicesCommonMigrationTaskServiceStateMigrationOptionsItemsEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["CONTINUOUS","DELETE_AFTER","USE_TRANSFORM_REQUEST","ALL_VERSIONS","ESTIMATE_COUNT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		comVmwareXenonServicesCommonMigrationTaskServiceStateMigrationOptionsItemsEnum = append(comVmwareXenonServicesCommonMigrationTaskServiceStateMigrationOptionsItemsEnum, v)
	}
}

func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) validateMigrationOptionsItemsEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, comVmwareXenonServicesCommonMigrationTaskServiceStateMigrationOptionsItemsEnum); err != nil {
		return err
	}
	return nil
}

func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) validateMigrationOptions(formats strfmt.Registry) error {

	if swag.IsZero(m.MigrationOptions) { // not required
		return nil
	}

	for i := 0; i < len(m.MigrationOptions); i++ {

		// value enum
		if err := m.validateMigrationOptionsItemsEnum("migrationOptions"+"."+strconv.Itoa(i), "body", m.MigrationOptions[i]); err != nil {
			return err
		}

	}

	return nil
}

func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) validateQuerySpec(formats strfmt.Registry) error {

	if swag.IsZero(m.QuerySpec) { // not required
		return nil
	}

	if m.QuerySpec != nil {

		if err := m.QuerySpec.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("querySpec")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) validateTaskInfo(formats strfmt.Registry) error {

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

// MarshalBinary interface implementation
func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareXenonServicesCommonMigrationTaskServiceState) UnmarshalBinary(b []byte) error {
	var res ComVmwareXenonServicesCommonMigrationTaskServiceState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
