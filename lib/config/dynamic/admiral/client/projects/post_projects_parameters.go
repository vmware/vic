package projects

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

	"github.com/vmware/vic/lib/config/dynamic/admiral/models"
)

// NewPostProjectsParams creates a new PostProjectsParams object
// with the default values initialized.
func NewPostProjectsParams() *PostProjectsParams {
	var ()
	return &PostProjectsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostProjectsParamsWithTimeout creates a new PostProjectsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostProjectsParamsWithTimeout(timeout time.Duration) *PostProjectsParams {
	var ()
	return &PostProjectsParams{

		timeout: timeout,
	}
}

// NewPostProjectsParamsWithContext creates a new PostProjectsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostProjectsParamsWithContext(ctx context.Context) *PostProjectsParams {
	var ()
	return &PostProjectsParams{

		Context: ctx,
	}
}

// NewPostProjectsParamsWithHTTPClient creates a new PostProjectsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostProjectsParamsWithHTTPClient(client *http.Client) *PostProjectsParams {
	var ()
	return &PostProjectsParams{
		HTTPClient: client,
	}
}

/*PostProjectsParams contains all the parameters to send to the API endpoint
for the post projects operation typically these are written to a http.Request
*/
type PostProjectsParams struct {

	/*Body*/
	Body *models.ComVmwareAdmiralAuthProjectProjectServiceProjectState

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post projects params
func (o *PostProjectsParams) WithTimeout(timeout time.Duration) *PostProjectsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post projects params
func (o *PostProjectsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post projects params
func (o *PostProjectsParams) WithContext(ctx context.Context) *PostProjectsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post projects params
func (o *PostProjectsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post projects params
func (o *PostProjectsParams) WithHTTPClient(client *http.Client) *PostProjectsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post projects params
func (o *PostProjectsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the post projects params
func (o *PostProjectsParams) WithBody(body *models.ComVmwareAdmiralAuthProjectProjectServiceProjectState) *PostProjectsParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the post projects params
func (o *PostProjectsParams) SetBody(body *models.ComVmwareAdmiralAuthProjectProjectServiceProjectState) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *PostProjectsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body == nil {
		o.Body = new(models.ComVmwareAdmiralAuthProjectProjectServiceProjectState)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
