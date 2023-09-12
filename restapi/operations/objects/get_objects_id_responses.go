// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// GetObjectsIDOKCode is the HTTP code returned for type GetObjectsIDOK
const GetObjectsIDOKCode int = 200

/*GetObjectsIDOK Object full info

swagger:response getObjectsIdOK
*/
type GetObjectsIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.AObject `json:"body,omitempty"`
}

// NewGetObjectsIDOK creates GetObjectsIDOK with default headers values
func NewGetObjectsIDOK() *GetObjectsIDOK {

	return &GetObjectsIDOK{}
}

// WithPayload adds the payload to the get objects Id o k response
func (o *GetObjectsIDOK) WithPayload(payload *models.AObject) *GetObjectsIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get objects Id o k response
func (o *GetObjectsIDOK) SetPayload(payload *models.AObject) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetObjectsIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetObjectsIDForbiddenCode is the HTTP code returned for type GetObjectsIDForbidden
const GetObjectsIDForbiddenCode int = 403

/*GetObjectsIDForbidden Forbidden

swagger:response getObjectsIdForbidden
*/
type GetObjectsIDForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetObjectsIDForbidden creates GetObjectsIDForbidden with default headers values
func NewGetObjectsIDForbidden() *GetObjectsIDForbidden {

	return &GetObjectsIDForbidden{}
}

// WithPayload adds the payload to the get objects Id forbidden response
func (o *GetObjectsIDForbidden) WithPayload(payload *models.Error) *GetObjectsIDForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get objects Id forbidden response
func (o *GetObjectsIDForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetObjectsIDForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetObjectsIDNotFoundCode is the HTTP code returned for type GetObjectsIDNotFound
const GetObjectsIDNotFoundCode int = 404

/*GetObjectsIDNotFound cant find object

swagger:response getObjectsIdNotFound
*/
type GetObjectsIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetObjectsIDNotFound creates GetObjectsIDNotFound with default headers values
func NewGetObjectsIDNotFound() *GetObjectsIDNotFound {

	return &GetObjectsIDNotFound{}
}

// WithPayload adds the payload to the get objects Id not found response
func (o *GetObjectsIDNotFound) WithPayload(payload *models.Error) *GetObjectsIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get objects Id not found response
func (o *GetObjectsIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetObjectsIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetObjectsIDDefault Unexpected error

swagger:response getObjectsIdDefault
*/
type GetObjectsIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetObjectsIDDefault creates GetObjectsIDDefault with default headers values
func NewGetObjectsIDDefault(code int) *GetObjectsIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetObjectsIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get objects ID default response
func (o *GetObjectsIDDefault) WithStatusCode(code int) *GetObjectsIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get objects ID default response
func (o *GetObjectsIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get objects ID default response
func (o *GetObjectsIDDefault) WithPayload(payload *models.Error) *GetObjectsIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get objects ID default response
func (o *GetObjectsIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetObjectsIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
