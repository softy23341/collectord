// Code generated by go-swagger; DO NOT EDIT.

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PutTaskIDOKCode is the HTTP code returned for type PutTaskIDOK
const PutTaskIDOKCode int = 200

/*PutTaskIDOK success

swagger:response putTaskIdOK
*/
type PutTaskIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.AGetTask `json:"body,omitempty"`
}

// NewPutTaskIDOK creates PutTaskIDOK with default headers values
func NewPutTaskIDOK() *PutTaskIDOK {

	return &PutTaskIDOK{}
}

// WithPayload adds the payload to the put task Id o k response
func (o *PutTaskIDOK) WithPayload(payload *models.AGetTask) *PutTaskIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put task Id o k response
func (o *PutTaskIDOK) SetPayload(payload *models.AGetTask) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTaskIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutTaskIDForbiddenCode is the HTTP code returned for type PutTaskIDForbidden
const PutTaskIDForbiddenCode int = 403

/*PutTaskIDForbidden Forbidden

swagger:response putTaskIdForbidden
*/
type PutTaskIDForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutTaskIDForbidden creates PutTaskIDForbidden with default headers values
func NewPutTaskIDForbidden() *PutTaskIDForbidden {

	return &PutTaskIDForbidden{}
}

// WithPayload adds the payload to the put task Id forbidden response
func (o *PutTaskIDForbidden) WithPayload(payload *models.Error) *PutTaskIDForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put task Id forbidden response
func (o *PutTaskIDForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTaskIDForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutTaskIDNotFoundCode is the HTTP code returned for type PutTaskIDNotFound
const PutTaskIDNotFoundCode int = 404

/*PutTaskIDNotFound cant find the task

swagger:response putTaskIdNotFound
*/
type PutTaskIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutTaskIDNotFound creates PutTaskIDNotFound with default headers values
func NewPutTaskIDNotFound() *PutTaskIDNotFound {

	return &PutTaskIDNotFound{}
}

// WithPayload adds the payload to the put task Id not found response
func (o *PutTaskIDNotFound) WithPayload(payload *models.Error) *PutTaskIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put task Id not found response
func (o *PutTaskIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTaskIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutTaskIDUnprocessableEntityCode is the HTTP code returned for type PutTaskIDUnprocessableEntity
const PutTaskIDUnprocessableEntityCode int = 422

/*PutTaskIDUnprocessableEntity validation error

swagger:response putTaskIdUnprocessableEntity
*/
type PutTaskIDUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutTaskIDUnprocessableEntity creates PutTaskIDUnprocessableEntity with default headers values
func NewPutTaskIDUnprocessableEntity() *PutTaskIDUnprocessableEntity {

	return &PutTaskIDUnprocessableEntity{}
}

// WithPayload adds the payload to the put task Id unprocessable entity response
func (o *PutTaskIDUnprocessableEntity) WithPayload(payload *models.Error) *PutTaskIDUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put task Id unprocessable entity response
func (o *PutTaskIDUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTaskIDUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PutTaskIDDefault Unexpected error

swagger:response putTaskIdDefault
*/
type PutTaskIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutTaskIDDefault creates PutTaskIDDefault with default headers values
func NewPutTaskIDDefault(code int) *PutTaskIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutTaskIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put task ID default response
func (o *PutTaskIDDefault) WithStatusCode(code int) *PutTaskIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put task ID default response
func (o *PutTaskIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put task ID default response
func (o *PutTaskIDDefault) WithPayload(payload *models.Error) *PutTaskIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put task ID default response
func (o *PutTaskIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutTaskIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
