// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	models "git.softndit.com/collector/backend/models"
)

// NewPostObjectsSearchParams creates a new PostObjectsSearchParams object
// no default values defined in spec.
func NewPostObjectsSearchParams() PostObjectsSearchParams {

	return PostObjectsSearchParams{}
}

// PostObjectsSearchParams contains all the bound params for the post objects search operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostObjectsSearch
type PostObjectsSearchParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	RSearchObjects *models.RSearchObjects
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostObjectsSearchParams() beforehand.
func (o *PostObjectsSearchParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	var errInterface interface{}

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RSearchObjects
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("rSearchObjects", "body", errInterface))
			} else {
				res = append(res, errors.NewParseError("rSearchObjects", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.RSearchObjects = &body
			}
		}
	} else {
		res = append(res, errors.Required("rSearchObjects", "body", errInterface))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
