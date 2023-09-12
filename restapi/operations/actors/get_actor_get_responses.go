// Code generated by go-swagger; DO NOT EDIT.

package actors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// GetActorGetOKCode is the HTTP code returned for type GetActorGetOK
const GetActorGetOKCode int = 200

/*GetActorGetOK actors list for selected root

swagger:response getActorGetOK
*/
type GetActorGetOK struct {

	/*
	  In: Body
	*/
	Payload *models.AActors `json:"body,omitempty"`
}

// NewGetActorGetOK creates GetActorGetOK with default headers values
func NewGetActorGetOK() *GetActorGetOK {

	return &GetActorGetOK{}
}

// WithPayload adds the payload to the get actor get o k response
func (o *GetActorGetOK) WithPayload(payload *models.AActors) *GetActorGetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get actor get o k response
func (o *GetActorGetOK) SetPayload(payload *models.AActors) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActorGetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetActorGetForbiddenCode is the HTTP code returned for type GetActorGetForbidden
const GetActorGetForbiddenCode int = 403

/*GetActorGetForbidden Forbidden

swagger:response getActorGetForbidden
*/
type GetActorGetForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetActorGetForbidden creates GetActorGetForbidden with default headers values
func NewGetActorGetForbidden() *GetActorGetForbidden {

	return &GetActorGetForbidden{}
}

// WithPayload adds the payload to the get actor get forbidden response
func (o *GetActorGetForbidden) WithPayload(payload *models.Error) *GetActorGetForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get actor get forbidden response
func (o *GetActorGetForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActorGetForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetActorGetDefault Unexpected error

swagger:response getActorGetDefault
*/
type GetActorGetDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetActorGetDefault creates GetActorGetDefault with default headers values
func NewGetActorGetDefault(code int) *GetActorGetDefault {
	if code <= 0 {
		code = 500
	}

	return &GetActorGetDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get actor get default response
func (o *GetActorGetDefault) WithStatusCode(code int) *GetActorGetDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get actor get default response
func (o *GetActorGetDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get actor get default response
func (o *GetActorGetDefault) WithPayload(payload *models.Error) *GetActorGetDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get actor get default response
func (o *GetActorGetDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetActorGetDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}