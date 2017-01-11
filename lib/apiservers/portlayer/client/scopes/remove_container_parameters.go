package scopes

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

// NewRemoveContainerParams creates a new RemoveContainerParams object
// with the default values initialized.
func NewRemoveContainerParams() *RemoveContainerParams {
	var ()
	return &RemoveContainerParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveContainerParamsWithTimeout creates a new RemoveContainerParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRemoveContainerParamsWithTimeout(timeout time.Duration) *RemoveContainerParams {
	var ()
	return &RemoveContainerParams{

		timeout: timeout,
	}
}

// NewRemoveContainerParamsWithContext creates a new RemoveContainerParams object
// with the default values initialized, and the ability to set a context for a request
func NewRemoveContainerParamsWithContext(ctx context.Context) *RemoveContainerParams {
	var ()
	return &RemoveContainerParams{

		Context: ctx,
	}
}

/*RemoveContainerParams contains all the parameters to send to the API endpoint
for the remove container operation typically these are written to a http.Request
*/
type RemoveContainerParams struct {

	/*Handle*/
	Handle string
	/*Scope*/
	Scope string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the remove container params
func (o *RemoveContainerParams) WithTimeout(timeout time.Duration) *RemoveContainerParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove container params
func (o *RemoveContainerParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove container params
func (o *RemoveContainerParams) WithContext(ctx context.Context) *RemoveContainerParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove container params
func (o *RemoveContainerParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHandle adds the handle to the remove container params
func (o *RemoveContainerParams) WithHandle(handle string) *RemoveContainerParams {
	o.SetHandle(handle)
	return o
}

// SetHandle adds the handle to the remove container params
func (o *RemoveContainerParams) SetHandle(handle string) {
	o.Handle = handle
}

// WithScope adds the scope to the remove container params
func (o *RemoveContainerParams) WithScope(scope string) *RemoveContainerParams {
	o.SetScope(scope)
	return o
}

// SetScope adds the scope to the remove container params
func (o *RemoveContainerParams) SetScope(scope string) {
	o.Scope = scope
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveContainerParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	// path param handle
	if err := r.SetPathParam("handle", o.Handle); err != nil {
		return err
	}

	// path param scope
	if err := r.SetPathParam("scope", o.Scope); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
