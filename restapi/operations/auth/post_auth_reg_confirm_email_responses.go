// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostAuthRegConfirmEmailNoContentCode is the HTTP code returned for type PostAuthRegConfirmEmailNoContent
const PostAuthRegConfirmEmailNoContentCode int = 204

/*PostAuthRegConfirmEmailNoContent Cinfirm email

swagger:response postAuthRegConfirmEmailNoContent
*/
type PostAuthRegConfirmEmailNoContent struct {
}

// NewPostAuthRegConfirmEmailNoContent creates PostAuthRegConfirmEmailNoContent with default headers values
func NewPostAuthRegConfirmEmailNoContent() *PostAuthRegConfirmEmailNoContent {

	return &PostAuthRegConfirmEmailNoContent{}
}

// WriteResponse to the client
func (o *PostAuthRegConfirmEmailNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostAuthRegConfirmEmailConflictCode is the HTTP code returned for type PostAuthRegConfirmEmailConflict
const PostAuthRegConfirmEmailConflictCode int = 409

/*PostAuthRegConfirmEmailConflict User already present

swagger:response postAuthRegConfirmEmailConflict
*/
type PostAuthRegConfirmEmailConflict struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostAuthRegConfirmEmailConflict creates PostAuthRegConfirmEmailConflict with default headers values
func NewPostAuthRegConfirmEmailConflict() *PostAuthRegConfirmEmailConflict {

	return &PostAuthRegConfirmEmailConflict{}
}

// WithPayload adds the payload to the post auth reg confirm email conflict response
func (o *PostAuthRegConfirmEmailConflict) WithPayload(payload *models.Error) *PostAuthRegConfirmEmailConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post auth reg confirm email conflict response
func (o *PostAuthRegConfirmEmailConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAuthRegConfirmEmailConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostAuthRegConfirmEmailDefault Unexpected error

swagger:response postAuthRegConfirmEmailDefault
*/
type PostAuthRegConfirmEmailDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostAuthRegConfirmEmailDefault creates PostAuthRegConfirmEmailDefault with default headers values
func NewPostAuthRegConfirmEmailDefault(code int) *PostAuthRegConfirmEmailDefault {
	if code <= 0 {
		code = 500
	}

	return &PostAuthRegConfirmEmailDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post auth reg confirm email default response
func (o *PostAuthRegConfirmEmailDefault) WithStatusCode(code int) *PostAuthRegConfirmEmailDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post auth reg confirm email default response
func (o *PostAuthRegConfirmEmailDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post auth reg confirm email default response
func (o *PostAuthRegConfirmEmailDefault) WithPayload(payload *models.Error) *PostAuthRegConfirmEmailDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post auth reg confirm email default response
func (o *PostAuthRegConfirmEmailDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAuthRegConfirmEmailDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}