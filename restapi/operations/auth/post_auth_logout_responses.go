// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostAuthLogoutNoContentCode is the HTTP code returned for type PostAuthLogoutNoContent
const PostAuthLogoutNoContentCode int = 204

/*PostAuthLogoutNoContent logout

swagger:response postAuthLogoutNoContent
*/
type PostAuthLogoutNoContent struct {
}

// NewPostAuthLogoutNoContent creates PostAuthLogoutNoContent with default headers values
func NewPostAuthLogoutNoContent() *PostAuthLogoutNoContent {

	return &PostAuthLogoutNoContent{}
}

// WriteResponse to the client
func (o *PostAuthLogoutNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*PostAuthLogoutDefault Unexpected error

swagger:response postAuthLogoutDefault
*/
type PostAuthLogoutDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostAuthLogoutDefault creates PostAuthLogoutDefault with default headers values
func NewPostAuthLogoutDefault(code int) *PostAuthLogoutDefault {
	if code <= 0 {
		code = 500
	}

	return &PostAuthLogoutDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post auth logout default response
func (o *PostAuthLogoutDefault) WithStatusCode(code int) *PostAuthLogoutDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post auth logout default response
func (o *PostAuthLogoutDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post auth logout default response
func (o *PostAuthLogoutDefault) WithPayload(payload *models.Error) *PostAuthLogoutDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post auth logout default response
func (o *PostAuthLogoutDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAuthLogoutDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}