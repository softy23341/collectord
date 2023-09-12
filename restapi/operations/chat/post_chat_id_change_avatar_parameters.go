// Code generated by go-swagger; DO NOT EDIT.

package chat

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"

	models "git.softndit.com/collector/backend/models"
)

// NewPostChatIDChangeAvatarParams creates a new PostChatIDChangeAvatarParams object
// no default values defined in spec.
func NewPostChatIDChangeAvatarParams() PostChatIDChangeAvatarParams {

	return PostChatIDChangeAvatarParams{}
}

// PostChatIDChangeAvatarParams contains all the bound params for the post chat ID change avatar operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostChatIDChangeAvatar
type PostChatIDChangeAvatarParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ID int64
	/*
	  Required: true
	  In: body
	*/
	REditChatAvatar *models.REditChatAvatar
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostChatIDChangeAvatarParams() beforehand.
func (o *PostChatIDChangeAvatarParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	var errInterface interface{}

	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("ID")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.REditChatAvatar
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("rEditChatAvatar", "body", errInterface))
			} else {
				res = append(res, errors.NewParseError("rEditChatAvatar", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.REditChatAvatar = &body
			}
		}
	} else {
		res = append(res, errors.Required("rEditChatAvatar", "body", errInterface))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostChatIDChangeAvatarParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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