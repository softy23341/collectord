// Code generated by go-swagger; DO NOT EDIT.

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// DeleteTaskIDNoContentCode is the HTTP code returned for type DeleteTaskIDNoContent
const DeleteTaskIDNoContentCode int = 204

/*DeleteTaskIDNoContent success

swagger:response deleteTaskIdNoContent
*/
type DeleteTaskIDNoContent struct {
}

// NewDeleteTaskIDNoContent creates DeleteTaskIDNoContent with default headers values
func NewDeleteTaskIDNoContent() *DeleteTaskIDNoContent {

	return &DeleteTaskIDNoContent{}
}

// WriteResponse to the client
func (o *DeleteTaskIDNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteTaskIDForbiddenCode is the HTTP code returned for type DeleteTaskIDForbidden
const DeleteTaskIDForbiddenCode int = 403

/*DeleteTaskIDForbidden Forbidden

swagger:response deleteTaskIdForbidden
*/
type DeleteTaskIDForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTaskIDForbidden creates DeleteTaskIDForbidden with default headers values
func NewDeleteTaskIDForbidden() *DeleteTaskIDForbidden {

	return &DeleteTaskIDForbidden{}
}

// WithPayload adds the payload to the delete task Id forbidden response
func (o *DeleteTaskIDForbidden) WithPayload(payload *models.Error) *DeleteTaskIDForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete task Id forbidden response
func (o *DeleteTaskIDForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTaskIDForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteTaskIDNotFoundCode is the HTTP code returned for type DeleteTaskIDNotFound
const DeleteTaskIDNotFoundCode int = 404

/*DeleteTaskIDNotFound cant find the task

swagger:response deleteTaskIdNotFound
*/
type DeleteTaskIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTaskIDNotFound creates DeleteTaskIDNotFound with default headers values
func NewDeleteTaskIDNotFound() *DeleteTaskIDNotFound {

	return &DeleteTaskIDNotFound{}
}

// WithPayload adds the payload to the delete task Id not found response
func (o *DeleteTaskIDNotFound) WithPayload(payload *models.Error) *DeleteTaskIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete task Id not found response
func (o *DeleteTaskIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTaskIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteTaskIDDefault Unexpected error

swagger:response deleteTaskIdDefault
*/
type DeleteTaskIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTaskIDDefault creates DeleteTaskIDDefault with default headers values
func NewDeleteTaskIDDefault(code int) *DeleteTaskIDDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteTaskIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete task ID default response
func (o *DeleteTaskIDDefault) WithStatusCode(code int) *DeleteTaskIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete task ID default response
func (o *DeleteTaskIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete task ID default response
func (o *DeleteTaskIDDefault) WithPayload(payload *models.Error) *DeleteTaskIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete task ID default response
func (o *DeleteTaskIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTaskIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
