package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// TaskRemoveConfig task remove config
// swagger:model TaskRemoveConfig
type TaskRemoveConfig struct {

	// handle
	// Required: true
	Handle interface{} `json:"handle"`

	// id
	// Required: true
	ID string `json:"id"`
}

// Validate validates this task remove config
func (m *TaskRemoveConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHandle(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskRemoveConfig) validateHandle(formats strfmt.Registry) error {

	return nil
}

func (m *TaskRemoveConfig) validateID(formats strfmt.Registry) error {

	if err := validate.RequiredString("id", "body", string(m.ID)); err != nil {
		return err
	}

	return nil
}
