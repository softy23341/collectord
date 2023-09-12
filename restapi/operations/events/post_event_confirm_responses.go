// Code generated by go-swagger; DO NOT EDIT.

package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostEventConfirmNoContentCode is the HTTP code returned for type PostEventConfirmNoContent
const PostEventConfirmNoContentCode int = 204

/*PostEventConfirmNoContent success

swagger:response postEventConfirmNoContent
*/
type PostEventConfirmNoContent struct {
}

// NewPostEventConfirmNoContent creates PostEventConfirmNoContent with default headers values
func NewPostEventConfirmNoContent() *PostEventConfirmNoContent {

	return &PostEventConfirmNoContent{}
}

// WriteResponse to the client
func (o *PostEventConfirmNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

/*PostEventConfirmDefault Unexpected error

swagger:response postEventConfirmDefault
*/
type PostEventConfirmDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostEventConfirmDefault creates PostEventConfirmDefault with default headers values
func NewPostEventConfirmDefault(code int) *PostEventConfirmDefault {
	if code <= 0 {
		code = 500
	}

	return &PostEventConfirmDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post event confirm default response
func (o *PostEventConfirmDefault) WithStatusCode(code int) *PostEventConfirmDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post event confirm default response
func (o *PostEventConfirmDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post event confirm default response
func (o *PostEventConfirmDefault) WithPayload(payload *models.Error) *PostEventConfirmDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post event confirm default response
func (o *PostEventConfirmDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEventConfirmDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
