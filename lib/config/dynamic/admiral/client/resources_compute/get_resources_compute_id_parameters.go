package resources_compute

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

// NewGetResourcesComputeIDParams creates a new GetResourcesComputeIDParams object
// with the default values initialized.
func NewGetResourcesComputeIDParams() *GetResourcesComputeIDParams {
	var ()
	return &GetResourcesComputeIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetResourcesComputeIDParamsWithTimeout creates a new GetResourcesComputeIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetResourcesComputeIDParamsWithTimeout(timeout time.Duration) *GetResourcesComputeIDParams {
	var ()
	return &GetResourcesComputeIDParams{

		timeout: timeout,
	}
}

// NewGetResourcesComputeIDParamsWithContext creates a new GetResourcesComputeIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetResourcesComputeIDParamsWithContext(ctx context.Context) *GetResourcesComputeIDParams {
	var ()
	return &GetResourcesComputeIDParams{

		Context: ctx,
	}
}

// NewGetResourcesComputeIDParamsWithHTTPClient creates a new GetResourcesComputeIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetResourcesComputeIDParamsWithHTTPClient(client *http.Client) *GetResourcesComputeIDParams {
	var ()
	return &GetResourcesComputeIDParams{
		HTTPClient: client,
	}
}

/*GetResourcesComputeIDParams contains all the parameters to send to the API endpoint
for the get resources compute ID operation typically these are written to a http.Request
*/
type GetResourcesComputeIDParams struct {

	/*ID*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get resources compute ID params
func (o *GetResourcesComputeIDParams) WithTimeout(timeout time.Duration) *GetResourcesComputeIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get resources compute ID params
func (o *GetResourcesComputeIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get resources compute ID params
func (o *GetResourcesComputeIDParams) WithContext(ctx context.Context) *GetResourcesComputeIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get resources compute ID params
func (o *GetResourcesComputeIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get resources compute ID params
func (o *GetResourcesComputeIDParams) WithHTTPClient(client *http.Client) *GetResourcesComputeIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get resources compute ID params
func (o *GetResourcesComputeIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get resources compute ID params
func (o *GetResourcesComputeIDParams) WithID(id string) *GetResourcesComputeIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get resources compute ID params
func (o *GetResourcesComputeIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetResourcesComputeIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
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
