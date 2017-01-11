package containers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetContainerInfoParams creates a new GetContainerInfoParams object
// with the default values initialized.
func NewGetContainerInfoParams() *GetContainerInfoParams {
	var ()
	return &GetContainerInfoParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetContainerInfoParamsWithTimeout creates a new GetContainerInfoParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetContainerInfoParamsWithTimeout(timeout time.Duration) *GetContainerInfoParams {
	var ()
	return &GetContainerInfoParams{

		timeout: timeout,
	}
}

// NewGetContainerInfoParamsWithContext creates a new GetContainerInfoParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetContainerInfoParamsWithContext(ctx context.Context) *GetContainerInfoParams {
	var ()
	return &GetContainerInfoParams{

		Context: ctx,
	}
}

/*GetContainerInfoParams contains all the parameters to send to the API endpoint
for the get container info operation typically these are written to a http.Request
*/
type GetContainerInfoParams struct {

	/*ID*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get container info params
func (o *GetContainerInfoParams) WithTimeout(timeout time.Duration) *GetContainerInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get container info params
func (o *GetContainerInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get container info params
func (o *GetContainerInfoParams) WithContext(ctx context.Context) *GetContainerInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get container info params
func (o *GetContainerInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithID adds the id to the get container info params
func (o *GetContainerInfoParams) WithID(id string) *GetContainerInfoParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get container info params
func (o *GetContainerInfoParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetContainerInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
