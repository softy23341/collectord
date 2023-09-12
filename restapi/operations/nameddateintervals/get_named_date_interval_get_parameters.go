// Code generated by go-swagger; DO NOT EDIT.

package nameddateintervals

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetNamedDateIntervalGetParams creates a new GetNamedDateIntervalGetParams object
// no default values defined in spec.
func NewGetNamedDateIntervalGetParams() GetNamedDateIntervalGetParams {

	return GetNamedDateIntervalGetParams{}
}

// GetNamedDateIntervalGetParams contains all the bound params for the get named date interval get operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetNamedDateIntervalGet
type GetNamedDateIntervalGetParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: query
	*/
	RootID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetNamedDateIntervalGetParams() beforehand.
func (o *GetNamedDateIntervalGetParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qRootID, qhkRootID, _ := qs.GetOK("rootID")
	if err := o.bindRootID(qRootID, qhkRootID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetNamedDateIntervalGetParams) bindRootID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var errInterface interface{}
	if !hasKey {
		return errors.Required("rootID", "query", errInterface)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("rootID", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("rootID", "query", "int64", raw)
	}
	o.RootID = value

	return nil
}
