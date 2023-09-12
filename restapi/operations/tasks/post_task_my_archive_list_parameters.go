// Code generated by go-swagger; DO NOT EDIT.

package tasks

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

// NewPostTaskMyArchiveListParams creates a new PostTaskMyArchiveListParams object
// no default values defined in spec.
func NewPostTaskMyArchiveListParams() PostTaskMyArchiveListParams {

	return PostTaskMyArchiveListParams{}
}

// PostTaskMyArchiveListParams contains all the bound params for the post task my archive list operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostTaskMyArchiveList
type PostTaskMyArchiveListParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	RMyArchiveList *models.RMyArchiveList
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostTaskMyArchiveListParams() beforehand.
func (o *PostTaskMyArchiveListParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	var errInterface interface{}

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RMyArchiveList
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("rMyArchiveList", "body", errInterface))
			} else {
				res = append(res, errors.NewParseError("rMyArchiveList", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.RMyArchiveList = &body
			}
		}
	} else {
		res = append(res, errors.Required("rMyArchiveList", "body", errInterface))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}