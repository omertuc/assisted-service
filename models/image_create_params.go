// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ImageCreateParams image create params
//
// swagger:model image-create-params
type ImageCreateParams struct {

	// SSH public key for debugging the installation.
	SSHPublicKey string `json:"ssh_public_key,omitempty"`

	// static ips config
	StaticIpsConfig []*StaticIPConfig `json:"static_ips_config"`
}

// Validate validates this image create params
func (m *ImageCreateParams) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStaticIpsConfig(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImageCreateParams) validateStaticIpsConfig(formats strfmt.Registry) error {

	if swag.IsZero(m.StaticIpsConfig) { // not required
		return nil
	}

	for i := 0; i < len(m.StaticIpsConfig); i++ {
		if swag.IsZero(m.StaticIpsConfig[i]) { // not required
			continue
		}

		if m.StaticIpsConfig[i] != nil {
			if err := m.StaticIpsConfig[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("static_ips_config" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ImageCreateParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ImageCreateParams) UnmarshalBinary(b []byte) error {
	var res ImageCreateParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
