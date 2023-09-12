// Code generated by go-swagger; DO NOT EDIT.

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostTaskOKCode is the HTTP code returned for type PostTaskOK
const PostTaskOKCode int = 200

/*PostTaskOK success

swagger:response postTaskOK
*/
type PostTaskOK struct {

	/*
	  In: Body
	*/
	Payload *models.AGetTask `json:"body,omitempty"`
}

// NewPostTaskOK creates PostTaskOK with default headers values
func NewPostTaskOK() *PostTaskOK {

	return &PostTaskOK{}
}

// WithPayload adds the payload to the post task o k response
func (o *PostTaskOK) WithPayload(payload *models.AGetTask) *PostTaskOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post task o k response
func (o *PostTaskOK) SetPayload(payload *models.AGetTask) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTaskOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostTaskForbiddenCode is the HTTP code returned for type PostTaskForbidden
const PostTaskForbiddenCode int = 403

/*PostTaskForbidden Forbidden

swagger:response postTaskForbidden
*/
type PostTaskForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostTaskForbidden creates PostTaskForbidden with default headers values
func NewPostTaskForbidden() *PostTaskForbidden {

	return &PostTaskForbidden{}
}

// WithPayload adds the payload to the post task forbidden response
func (o *PostTaskForbidden) WithPayload(payload *models.Error) *PostTaskForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post task forbidden response
func (o *PostTaskForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTaskForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostTaskUnprocessableEntityCode is the HTTP code returned for type PostTaskUnprocessableEntity
const PostTaskUnprocessableEntityCode int = 422

/*PostTaskUnprocessableEntity validation error

swagger:response postTaskUnprocessableEntity
*/
type PostTaskUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostTaskUnprocessableEntity creates PostTaskUnprocessableEntity with default headers values
func NewPostTaskUnprocessableEntity() *PostTaskUnprocessableEntity {

	return &PostTaskUnprocessableEntity{}
}

// WithPayload adds the payload to the post task unprocessable entity response
func (o *PostTaskUnprocessableEntity) WithPayload(payload *models.Error) *PostTaskUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post task unprocessable entity response
func (o *PostTaskUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTaskUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostTaskDefault Unexpected error

swagger:response postTaskDefault
*/
type PostTaskDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostTaskDefault creates PostTaskDefault with default headers values
func NewPostTaskDefault(code int) *PostTaskDefault {
	if code <= 0 {
		code = 500
	}

	return &PostTaskDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post task default response
func (o *PostTaskDefault) WithStatusCode(code int) *PostTaskDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post task default response
func (o *PostTaskDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post task default response
func (o *PostTaskDefault) WithPayload(payload *models.Error) *PostTaskDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post task default response
func (o *PostTaskDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTaskDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}