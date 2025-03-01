// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GlowReminder glow reminder
//
// swagger:model GlowReminder
type GlowReminder struct {

	// colour
	Colour int64 `json:"colour,omitempty"`

	// mode
	Mode int64 `json:"mode,omitempty"`
}

// Validate validates this glow reminder
func (m *GlowReminder) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this glow reminder based on context it is used
func (m *GlowReminder) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GlowReminder) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GlowReminder) UnmarshalBinary(b []byte) error {
	var res GlowReminder
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
