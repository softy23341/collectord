// Code generated by go-swagger; DO NOT EDIT.

package originlocations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteOriginLocationIDDeleteParams creates a new DeleteOriginLocationIDDeleteParams object
// no default values defined in spec.
func NewDeleteOriginLocationIDDeleteParams() DeleteOriginLocationIDDeleteParams {

	return DeleteOriginLocationIDDeleteParams{}
}

// DeleteOriginLocationIDDeleteParams contains all the bound params for the delete origin location ID delete operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteOriginLocationIDDelete
type DeleteOriginLocationIDDeleteParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteOriginLocationIDDeleteParams() beforehand.
func (o *DeleteOriginLocationIDDeleteParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("ID")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeleteOriginLocationIDDeleteParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("ID", "path", "int64", raw)
	}
	o.ID = value

	return nil
}
