package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState com vmware photon controller model tasks monitoring stats aggregation task service stats aggregation task state
// swagger:model com:vmware:photon:controller:model:tasks:monitoring:StatsAggregationTaskService:StatsAggregationTaskState
type ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState struct {

	// failure message
	FailureMessage string `json:"failureMessage,omitempty"`

	// The set of metric names to aggregate on
	MetricNames []string `json:"metricNames"`

	// The query to lookup resources for stats aggregation
	Query *ComVmwareXenonServicesCommonQueryTaskQuery `json:"query,omitempty"`

	// task info
	TaskInfo *ComVmwareXenonCommonTaskState `json:"taskInfo,omitempty"`
}

// Validate validates this com vmware photon controller model tasks monitoring stats aggregation task service stats aggregation task state
func (m *ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetricNames(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateQuery(formats); err != nil {
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

func (m *ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState) validateMetricNames(formats strfmt.Registry) error {

	if swag.IsZero(m.MetricNames) { // not required
		return nil
	}

	return nil
}

func (m *ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState) validateQuery(formats strfmt.Registry) error {

	if swag.IsZero(m.Query) { // not required
		return nil
	}

	if m.Query != nil {

		if err := m.Query.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("query")
			}
			return err
		}
	}

	return nil
}

func (m *ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState) validateTaskInfo(formats strfmt.Registry) error {

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
func (m *ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState) UnmarshalBinary(b []byte) error {
	var res ComVmwarePhotonControllerModelTasksMonitoringStatsAggregationTaskServiceStatsAggregationTaskState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
