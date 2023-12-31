// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	models "git.softndit.com/collector/backend/models"
)

// NewPostAuthPasswordResetTokenParams creates a new PostAuthPasswordResetTokenParams object
// no default values defined in spec.
func NewPostAuthPasswordResetTokenParams() PostAuthPasswordResetTokenParams {

	return PostAuthPasswordResetTokenParams{}
}

// PostAuthPasswordResetTokenParams contains all the bound params for the post auth password reset token operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostAuthPasswordResetToken
type PostAuthPasswordResetTokenParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	RPasswordReset *models.RPasswordReset
	/*
	  Required: true
	  In: path
	*/
	Token string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostAuthPasswordResetTokenParams() beforehand.
func (o *PostAuthPasswordResetTokenParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	var errInterface interface{}

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RPasswordReset
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("rPasswordReset", "body", errInterface))
			} else {
				res = append(res, errors.NewParseError("rPasswordReset", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.RPasswordReset = &body
			}
		}
	} else {
		res = append(res, errors.Required("rPasswordReset", "body", errInterface))
	}
	rToken, rhkToken, _ := route.Params.GetOK("token")
	if err := o.bindToken(rToken, rhkToken, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostAuthPasswordResetTokenParams) bindToken(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Token = raw

	return nil
}
