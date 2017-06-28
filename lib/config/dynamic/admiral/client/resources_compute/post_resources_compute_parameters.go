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

	"github.com/vmware/vic/lib/config/dynamic/admiral/models"
)

// NewPostResourcesComputeParams creates a new PostResourcesComputeParams object
// with the default values initialized.
func NewPostResourcesComputeParams() *PostResourcesComputeParams {
	var ()
	return &PostResourcesComputeParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostResourcesComputeParamsWithTimeout creates a new PostResourcesComputeParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostResourcesComputeParamsWithTimeout(timeout time.Duration) *PostResourcesComputeParams {
	var ()
	return &PostResourcesComputeParams{

		timeout: timeout,
	}
}

// NewPostResourcesComputeParamsWithContext creates a new PostResourcesComputeParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostResourcesComputeParamsWithContext(ctx context.Context) *PostResourcesComputeParams {
	var ()
	return &PostResourcesComputeParams{

		Context: ctx,
	}
}

// NewPostResourcesComputeParamsWithHTTPClient creates a new PostResourcesComputeParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostResourcesComputeParamsWithHTTPClient(client *http.Client) *PostResourcesComputeParams {
	var ()
	return &PostResourcesComputeParams{
		HTTPClient: client,
	}
}

/*PostResourcesComputeParams contains all the parameters to send to the API endpoint
for the post resources compute operation typically these are written to a http.Request
*/
type PostResourcesComputeParams struct {

	/*Body*/
	Body *models.ComVmwarePhotonControllerModelResourcesComputeServiceComputeState

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post resources compute params
func (o *PostResourcesComputeParams) WithTimeout(timeout time.Duration) *PostResourcesComputeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post resources compute params
func (o *PostResourcesComputeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post resources compute params
func (o *PostResourcesComputeParams) WithContext(ctx context.Context) *PostResourcesComputeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post resources compute params
func (o *PostResourcesComputeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post resources compute params
func (o *PostResourcesComputeParams) WithHTTPClient(client *http.Client) *PostResourcesComputeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post resources compute params
func (o *PostResourcesComputeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the post resources compute params
func (o *PostResourcesComputeParams) WithBody(body *models.ComVmwarePhotonControllerModelResourcesComputeServiceComputeState) *PostResourcesComputeParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the post resources compute params
func (o *PostResourcesComputeParams) SetBody(body *models.ComVmwarePhotonControllerModelResourcesComputeServiceComputeState) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *PostResourcesComputeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body == nil {
		o.Body = new(models.ComVmwarePhotonControllerModelResourcesComputeServiceComputeState)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
