// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PutObjectsIDNoContentCode is the HTTP code returned for type PutObjectsIDNoContent
const PutObjectsIDNoContentCode int = 204

/*PutObjectsIDNoContent ok

swagger:response putObjectsIdNoContent
*/
type PutObjectsIDNoContent struct {
}

// NewPutObjectsIDNoContent creates PutObjectsIDNoContent with default headers values
func NewPutObjectsIDNoContent() *PutObjectsIDNoContent {

	return &PutObjectsIDNoContent{}
}

// WriteResponse to the client
func (o *PutObjectsIDNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PutObjectsIDForbiddenCode is the HTTP code returned for type PutObjectsIDForbidden
const PutObjectsIDForbiddenCode int = 403

/*PutObjectsIDForbidden Forbidden

swagger:response putObjectsIdForbidden
*/
type PutObjectsIDForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutObjectsIDForbidden creates PutObjectsIDForbidden with default headers values
func NewPutObjectsIDForbidden() *PutObjectsIDForbidden {

	return &PutObjectsIDForbidden{}
}

// WithPayload adds the payload to the put objects Id forbidden response
func (o *PutObjectsIDForbidden) WithPayload(payload *models.Error) *PutObjectsIDForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put objects Id forbidden response
func (o *PutObjectsIDForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutObjectsIDForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutObjectsIDNotFoundCode is the HTTP code returned for type PutObjectsIDNotFound
const PutObjectsIDNotFoundCode int = 404

/*PutObjectsIDNotFound cant find object

swagger:response putObjectsIdNotFound
*/
type PutObjectsIDNotFound struct {
}

// NewPutObjectsIDNotFound creates PutObjectsIDNotFound with default headers values
func NewPutObjectsIDNotFound() *PutObjectsIDNotFound {

	return &PutObjectsIDNotFound{}
}

// WriteResponse to the client
func (o *PutObjectsIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

/*PutObjectsIDDefault Unexpected error

swagger:response putObjectsIdDefault
*/
type PutObjectsIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutObjectsIDDefault creates PutObjectsIDDefault with default headers values
func NewPutObjectsIDDefault(code int) *PutObjectsIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutObjectsIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put objects ID default response
func (o *PutObjectsIDDefault) WithStatusCode(code int) *PutObjectsIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put objects ID default response
func (o *PutObjectsIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put objects ID default response
func (o *PutObjectsIDDefault) WithPayload(payload *models.Error) *PutObjectsIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put objects ID default response
func (o *PutObjectsIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutObjectsIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
