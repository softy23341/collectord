// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostAuthPasswordRecoveryNoContentCode is the HTTP code returned for type PostAuthPasswordRecoveryNoContent
const PostAuthPasswordRecoveryNoContentCode int = 204

/*PostAuthPasswordRecoveryNoContent email was send

swagger:response postAuthPasswordRecoveryNoContent
*/
type PostAuthPasswordRecoveryNoContent struct {
}

// NewPostAuthPasswordRecoveryNoContent creates PostAuthPasswordRecoveryNoContent with default headers values
func NewPostAuthPasswordRecoveryNoContent() *PostAuthPasswordRecoveryNoContent {

	return &PostAuthPasswordRecoveryNoContent{}
}

// WriteResponse to the client
func (o *PostAuthPasswordRecoveryNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*PostAuthPasswordRecoveryDefault Unexpected error

swagger:response postAuthPasswordRecoveryDefault
*/
type PostAuthPasswordRecoveryDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostAuthPasswordRecoveryDefault creates PostAuthPasswordRecoveryDefault with default headers values
func NewPostAuthPasswordRecoveryDefault(code int) *PostAuthPasswordRecoveryDefault {
	if code <= 0 {
		code = 500
	}

	return &PostAuthPasswordRecoveryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post auth password recovery default response
func (o *PostAuthPasswordRecoveryDefault) WithStatusCode(code int) *PostAuthPasswordRecoveryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post auth password recovery default response
func (o *PostAuthPasswordRecoveryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post auth password recovery default response
func (o *PostAuthPasswordRecoveryDefault) WithPayload(payload *models.Error) *PostAuthPasswordRecoveryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post auth password recovery default response
func (o *PostAuthPasswordRecoveryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAuthPasswordRecoveryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
