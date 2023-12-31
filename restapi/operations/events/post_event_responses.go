// Code generated by go-swagger; DO NOT EDIT.

package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostEventOKCode is the HTTP code returned for type PostEventOK
const PostEventOKCode int = 200

/*PostEventOK Events list

swagger:response postEventOK
*/
type PostEventOK struct {

	/*
	  In: Body
	*/
	Payload *models.AGetEvents `json:"body,omitempty"`
}

// NewPostEventOK creates PostEventOK with default headers values
func NewPostEventOK() *PostEventOK {

	return &PostEventOK{}
}

// WithPayload adds the payload to the post event o k response
func (o *PostEventOK) WithPayload(payload *models.AGetEvents) *PostEventOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post event o k response
func (o *PostEventOK) SetPayload(payload *models.AGetEvents) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEventOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostEventDefault Unexpected error

swagger:response postEventDefault
*/
type PostEventDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostEventDefault creates PostEventDefault with default headers values
func NewPostEventDefault(code int) *PostEventDefault {
	if code <= 0 {
		code = 500
	}

	return &PostEventDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post event default response
func (o *PostEventDefault) WithStatusCode(code int) *PostEventDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post event default response
func (o *PostEventDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post event default response
func (o *PostEventDefault) WithPayload(payload *models.Error) *PostEventDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post event default response
func (o *PostEventDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEventDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
