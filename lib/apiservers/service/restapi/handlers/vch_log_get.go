package handlers

import (
	"net/url"

	"github.com/vmware/vic/lib/install/vchlog"

	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/vic/lib/apiservers/service/models"
	"github.com/vmware/vic/lib/apiservers/service/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/service/restapi/handlers/util"

)

type VCHLogGet struct {
}

type VCHDatacenterLogGet struct {
}

func (h *VCHLogGet) Handle(params operations.GetTargetTargetVchVchIDLogParams, principal interface{}) middleware.Responder {
	d, err := buildData(params.HTTPRequest.Context(),
		url.URL{Host: params.Target},
		principal.(Credentials).user,
		principal.(Credentials).pass,
		params.Thumbprint,
		nil,
		nil)
	if err != nil {
		return operations.NewGetTargetTargetVchVchIDLogDefault(util.StatusCode(err)).WithPayload(&models.Error{Message: err.Error()})
	}

	reader, err := vchlog.DownloadLogFile(params.HTTPRequest.Context())
	if err != nil {
		return operations.NewGetTargetTargetVchVchIDLogDefault()
	}


	return middleware.NotImplemented("Not Implemented")
}

func (h *VCHDatacenterLogGet) Handle(params operations.GetTargetTargetDatacenterDatacenterVchVchIDLogParams, principal interface{}) middleware.Responder {
	return middleware.NotImplemented("Not Implemented")
}