// Code generated by go-swagger; DO NOT EDIT.

package medias

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostMediasOKCode is the HTTP code returned for type PostMediasOK
const PostMediasOKCode int = 200

/*PostMediasOK Success

swagger:response postMediasOK
*/
type PostMediasOK struct {

	/*
	  In: Body
	*/
	Payload *models.ANewMedia `json:"body,omitempty"`
}

// NewPostMediasOK creates PostMediasOK with default headers values
func NewPostMediasOK() *PostMediasOK {

	return &PostMediasOK{}
}

// WithPayload adds the payload to the post medias o k response
func (o *PostMediasOK) WithPayload(payload *models.ANewMedia) *PostMediasOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post medias o k response
func (o *PostMediasOK) SetPayload(payload *models.ANewMedia) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMediasOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostMediasDefault Unexpected error

swagger:response postMediasDefault
*/
type PostMediasDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMediasDefault creates PostMediasDefault with default headers values
func NewPostMediasDefault(code int) *PostMediasDefault {
	if code <= 0 {
		code = 500
	}

	return &PostMediasDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post medias default response
func (o *PostMediasDefault) WithStatusCode(code int) *PostMediasDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post medias default response
func (o *PostMediasDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post medias default response
func (o *PostMediasDefault) WithPayload(payload *models.Error) *PostMediasDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post medias default response
func (o *PostMediasDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMediasDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
