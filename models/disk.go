// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Disk disk
//
// swagger:model disk
type Disk struct {

	// bootable
	Bootable bool `json:"bootable,omitempty"`

	// by path
	ByPath string `json:"by_path,omitempty"`

	// drive type
	DriveType string `json:"drive_type,omitempty"`

	// hctl
	Hctl string `json:"hctl,omitempty"`

	// installation eligibility
	InstallationEligibility *DiskInstallationEligibility `json:"installation_eligibility,omitempty"`

	// model
	Model string `json:"model,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// path
	Path string `json:"path,omitempty"`

	// serial
	Serial string `json:"serial,omitempty"`

	// size bytes
	SizeBytes int64 `json:"size_bytes,omitempty"`

	// smart
	Smart string `json:"smart,omitempty"`

	// vendor
	Vendor string `json:"vendor,omitempty"`

	// wwn
	Wwn string `json:"wwn,omitempty"`
}

// Validate validates this disk
func (m *Disk) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInstallationEligibility(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Disk) validateInstallationEligibility(formats strfmt.Registry) error {

	if swag.IsZero(m.InstallationEligibility) { // not required
		return nil
	}

	if m.InstallationEligibility != nil {
		if err := m.InstallationEligibility.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("installation_eligibility")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Disk) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Disk) UnmarshalBinary(b []byte) error {
	var res Disk
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DiskInstallationEligibility disk installation eligibility
//
// swagger:model DiskInstallationEligibility
type DiskInstallationEligibility struct {

	// Whether the disk is eligible for installation or not.
	Eligible bool `json:"eligible,omitempty"`

	// Reasons for why this disk is not elligible for installation.
	NotEligibleReasons []string `json:"not_eligible_reasons"`
}

// Validate validates this disk installation eligibility
func (m *DiskInstallationEligibility) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DiskInstallationEligibility) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DiskInstallationEligibility) UnmarshalBinary(b []byte) error {
	var res DiskInstallationEligibility
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
