// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// SubscriberState subscriber state
// swagger:model subscriber_state
type SubscriberState struct {

	// lte auth next seq
	LteAuthNextSeq uint64 `json:"lte_auth_next_seq,omitempty"`

	// tgpp aaa server name
	TgppAaaServerName string `json:"tgpp_aaa_server_name,omitempty"`

	// tgpp aaa server registered
	TgppAaaServerRegistered bool `json:"tgpp_aaa_server_registered,omitempty"`
}

// Validate validates this subscriber state
func (m *SubscriberState) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SubscriberState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriberState) UnmarshalBinary(b []byte) error {
	var res SubscriberState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
