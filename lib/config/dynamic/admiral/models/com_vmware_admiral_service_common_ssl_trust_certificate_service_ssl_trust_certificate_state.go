package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState com vmware admiral service common ssl trust certificate service ssl trust certificate state
// swagger:model com:vmware:admiral:service:common:SslTrustCertificateService:SslTrustCertificateState
type ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState struct {

	// The SSL trust certificate encoded into .PEM format.
	Certificate string `json:"certificate,omitempty"`

	// The common name of the certificate.
	CommonName string `json:"commonName,omitempty"`

	// The fingerprint of the certificate in SHA-1 form.
	Fingerprint string `json:"fingerprint,omitempty"`

	// The issuer name of the certificate.
	IssuerName string `json:"issuerName,omitempty"`

	// The serial of the certificate.
	Serial string `json:"serial,omitempty"`

	// tenant links
	TenantLinks []string `json:"tenantLinks"`

	// The date since the certificate is valid.
	ValidSince int64 `json:"validSince,omitempty"`

	// The date until the certificate is valid.
	ValidTo int64 `json:"validTo,omitempty"`
}

// Validate validates this com vmware admiral service common ssl trust certificate service ssl trust certificate state
func (m *ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTenantLinks(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState) validateTenantLinks(formats strfmt.Registry) error {

	if swag.IsZero(m.TenantLinks) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState) UnmarshalBinary(b []byte) error {
	var res ComVmwareAdmiralServiceCommonSslTrustCertificateServiceSslTrustCertificateState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
