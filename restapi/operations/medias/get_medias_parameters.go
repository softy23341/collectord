// Code generated by go-swagger; DO NOT EDIT.

package medias

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

// NewGetMediasParams creates a new GetMediasParams object
// no default values defined in spec.
func NewGetMediasParams() GetMediasParams {

	return GetMediasParams{}
}

// GetMediasParams contains all the bound params for the get medias operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetMedias
type GetMediasParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	MessageID *int64
	/*
	  Required: true
	  In: query
	*/
	Path string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetMediasParams() beforehand.
func (o *GetMediasParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qMessageID, qhkMessageID, _ := qs.GetOK("messageID")
	if err := o.bindMessageID(qMessageID, qhkMessageID, route.Formats); err != nil {
		res = append(res, err)
	}

	qPath, qhkPath, _ := qs.GetOK("path")
	if err := o.bindPath(qPath, qhkPath, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetMediasParams) bindMessageID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("messageID", "query", "int64", raw)
	}
	o.MessageID = &value

	return nil
}

func (o *GetMediasParams) bindPath(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var errInterface interface{}
	if !hasKey {
		return errors.Required("path", "query", errInterface)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("path", "query", raw); err != nil {
		return err
	}

	o.Path = raw

	return nil
}
