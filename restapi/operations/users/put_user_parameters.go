// Code generated by go-swagger; DO NOT EDIT.

package users

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

// NewPutUserParams creates a new PutUserParams object
// no default values defined in spec.
func NewPutUserParams() PutUserParams {

	return PutUserParams{}
}

// PutUserParams contains all the bound params for the put user operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutUser
type PutUserParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	REditUser *models.REditUser
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutUserParams() beforehand.
func (o *PutUserParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	var errInterface interface{}

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.REditUser
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("rEditUser", "body", errInterface))
			} else {
				res = append(res, errors.NewParseError("rEditUser", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.REditUser = &body
			}
		}
	} else {
		res = append(res, errors.Required("rEditUser", "body", errInterface))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
