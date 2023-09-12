// Code generated by go-swagger; DO NOT EDIT.

package nameddateintervals

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// DeleteNamedDateIntervalIDDeleteNoContentCode is the HTTP code returned for type DeleteNamedDateIntervalIDDeleteNoContent
const DeleteNamedDateIntervalIDDeleteNoContentCode int = 204

/*DeleteNamedDateIntervalIDDeleteNoContent succes

swagger:response deleteNamedDateIntervalIdDeleteNoContent
*/
type DeleteNamedDateIntervalIDDeleteNoContent struct {
}

// NewDeleteNamedDateIntervalIDDeleteNoContent creates DeleteNamedDateIntervalIDDeleteNoContent with default headers values
func NewDeleteNamedDateIntervalIDDeleteNoContent() *DeleteNamedDateIntervalIDDeleteNoContent {

	return &DeleteNamedDateIntervalIDDeleteNoContent{}
}

// WriteResponse to the client
func (o *DeleteNamedDateIntervalIDDeleteNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteNamedDateIntervalIDDeleteForbiddenCode is the HTTP code returned for type DeleteNamedDateIntervalIDDeleteForbidden
const DeleteNamedDateIntervalIDDeleteForbiddenCode int = 403

/*DeleteNamedDateIntervalIDDeleteForbidden access forbidden

swagger:response deleteNamedDateIntervalIdDeleteForbidden
*/
type DeleteNamedDateIntervalIDDeleteForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteNamedDateIntervalIDDeleteForbidden creates DeleteNamedDateIntervalIDDeleteForbidden with default headers values
func NewDeleteNamedDateIntervalIDDeleteForbidden() *DeleteNamedDateIntervalIDDeleteForbidden {

	return &DeleteNamedDateIntervalIDDeleteForbidden{}
}

// WithPayload adds the payload to the delete named date interval Id delete forbidden response
func (o *DeleteNamedDateIntervalIDDeleteForbidden) WithPayload(payload *models.Error) *DeleteNamedDateIntervalIDDeleteForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete named date interval Id delete forbidden response
func (o *DeleteNamedDateIntervalIDDeleteForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteNamedDateIntervalIDDeleteForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteNamedDateIntervalIDDeleteNotFoundCode is the HTTP code returned for type DeleteNamedDateIntervalIDDeleteNotFound
const DeleteNamedDateIntervalIDDeleteNotFoundCode int = 404

/*DeleteNamedDateIntervalIDDeleteNotFound cant find named date interval

swagger:response deleteNamedDateIntervalIdDeleteNotFound
*/
type DeleteNamedDateIntervalIDDeleteNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteNamedDateIntervalIDDeleteNotFound creates DeleteNamedDateIntervalIDDeleteNotFound with default headers values
func NewDeleteNamedDateIntervalIDDeleteNotFound() *DeleteNamedDateIntervalIDDeleteNotFound {

	return &DeleteNamedDateIntervalIDDeleteNotFound{}
}

// WithPayload adds the payload to the delete named date interval Id delete not found response
func (o *DeleteNamedDateIntervalIDDeleteNotFound) WithPayload(payload *models.Error) *DeleteNamedDateIntervalIDDeleteNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete named date interval Id delete not found response
func (o *DeleteNamedDateIntervalIDDeleteNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteNamedDateIntervalIDDeleteNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteNamedDateIntervalIDDeleteDefault Unexpected error

swagger:response deleteNamedDateIntervalIdDeleteDefault
*/
type DeleteNamedDateIntervalIDDeleteDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteNamedDateIntervalIDDeleteDefault creates DeleteNamedDateIntervalIDDeleteDefault with default headers values
func NewDeleteNamedDateIntervalIDDeleteDefault(code int) *DeleteNamedDateIntervalIDDeleteDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteNamedDateIntervalIDDeleteDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete named date interval ID delete default response
func (o *DeleteNamedDateIntervalIDDeleteDefault) WithStatusCode(code int) *DeleteNamedDateIntervalIDDeleteDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete named date interval ID delete default response
func (o *DeleteNamedDateIntervalIDDeleteDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete named date interval ID delete default response
func (o *DeleteNamedDateIntervalIDDeleteDefault) WithPayload(payload *models.Error) *DeleteNamedDateIntervalIDDeleteDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete named date interval ID delete default response
func (o *DeleteNamedDateIntervalIDDeleteDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteNamedDateIntervalIDDeleteDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
