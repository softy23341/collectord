// Code generated by go-swagger; DO NOT EDIT.

package collections

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// GetCollectionsDraftOKCode is the HTTP code returned for type GetCollectionsDraftOK
const GetCollectionsDraftOKCode int = 200

/*GetCollectionsDraftOK success

swagger:response getCollectionsDraftOK
*/
type GetCollectionsDraftOK struct {

	/*
	  In: Body
	*/
	Payload *models.Collection `json:"body,omitempty"`
}

// NewGetCollectionsDraftOK creates GetCollectionsDraftOK with default headers values
func NewGetCollectionsDraftOK() *GetCollectionsDraftOK {

	return &GetCollectionsDraftOK{}
}

// WithPayload adds the payload to the get collections draft o k response
func (o *GetCollectionsDraftOK) WithPayload(payload *models.Collection) *GetCollectionsDraftOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get collections draft o k response
func (o *GetCollectionsDraftOK) SetPayload(payload *models.Collection) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCollectionsDraftOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetCollectionsDraftDefault Unexpected error

swagger:response getCollectionsDraftDefault
*/
type GetCollectionsDraftDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetCollectionsDraftDefault creates GetCollectionsDraftDefault with default headers values
func NewGetCollectionsDraftDefault(code int) *GetCollectionsDraftDefault {
	if code <= 0 {
		code = 500
	}

	return &GetCollectionsDraftDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get collections draft default response
func (o *GetCollectionsDraftDefault) WithStatusCode(code int) *GetCollectionsDraftDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get collections draft default response
func (o *GetCollectionsDraftDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get collections draft default response
func (o *GetCollectionsDraftDefault) WithPayload(payload *models.Error) *GetCollectionsDraftDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get collections draft default response
func (o *GetCollectionsDraftDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCollectionsDraftDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
