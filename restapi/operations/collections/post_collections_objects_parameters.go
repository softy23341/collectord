// Code generated by go-swagger; DO NOT EDIT.

package collections

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

// NewPostCollectionsObjectsParams creates a new PostCollectionsObjectsParams object
// no default values defined in spec.
func NewPostCollectionsObjectsParams() PostCollectionsObjectsParams {

	return PostCollectionsObjectsParams{}
}

// PostCollectionsObjectsParams contains all the bound params for the post collections objects operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostCollectionsObjects
type PostCollectionsObjectsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	RGetCollectionsObjects *models.RGetCollectionsObjects
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostCollectionsObjectsParams() beforehand.
func (o *PostCollectionsObjectsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	var errInterface interface{}

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RGetCollectionsObjects
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("rGetCollectionsObjects", "body", errInterface))
			} else {
				res = append(res, errors.NewParseError("rGetCollectionsObjects", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.RGetCollectionsObjects = &body
			}
		}
	} else {
		res = append(res, errors.Required("rGetCollectionsObjects", "body", errInterface))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
