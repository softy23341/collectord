// Code generated by go-swagger; DO NOT EDIT.

package originlocations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// DeleteOriginLocationIDDeleteNoContentCode is the HTTP code returned for type DeleteOriginLocationIDDeleteNoContent
const DeleteOriginLocationIDDeleteNoContentCode int = 204

/*DeleteOriginLocationIDDeleteNoContent succes

swagger:response deleteOriginLocationIdDeleteNoContent
*/
type DeleteOriginLocationIDDeleteNoContent struct {
}

// NewDeleteOriginLocationIDDeleteNoContent creates DeleteOriginLocationIDDeleteNoContent with default headers values
func NewDeleteOriginLocationIDDeleteNoContent() *DeleteOriginLocationIDDeleteNoContent {

	return &DeleteOriginLocationIDDeleteNoContent{}
}

// WriteResponse to the client
func (o *DeleteOriginLocationIDDeleteNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteOriginLocationIDDeleteForbiddenCode is the HTTP code returned for type DeleteOriginLocationIDDeleteForbidden
const DeleteOriginLocationIDDeleteForbiddenCode int = 403

/*DeleteOriginLocationIDDeleteForbidden access forbidden

swagger:response deleteOriginLocationIdDeleteForbidden
*/
type DeleteOriginLocationIDDeleteForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteOriginLocationIDDeleteForbidden creates DeleteOriginLocationIDDeleteForbidden with default headers values
func NewDeleteOriginLocationIDDeleteForbidden() *DeleteOriginLocationIDDeleteForbidden {

	return &DeleteOriginLocationIDDeleteForbidden{}
}

// WithPayload adds the payload to the delete origin location Id delete forbidden response
func (o *DeleteOriginLocationIDDeleteForbidden) WithPayload(payload *models.Error) *DeleteOriginLocationIDDeleteForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete origin location Id delete forbidden response
func (o *DeleteOriginLocationIDDeleteForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteOriginLocationIDDeleteForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteOriginLocationIDDeleteNotFoundCode is the HTTP code returned for type DeleteOriginLocationIDDeleteNotFound
const DeleteOriginLocationIDDeleteNotFoundCode int = 404

/*DeleteOriginLocationIDDeleteNotFound cant find origin location

swagger:response deleteOriginLocationIdDeleteNotFound
*/
type DeleteOriginLocationIDDeleteNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteOriginLocationIDDeleteNotFound creates DeleteOriginLocationIDDeleteNotFound with default headers values
func NewDeleteOriginLocationIDDeleteNotFound() *DeleteOriginLocationIDDeleteNotFound {

	return &DeleteOriginLocationIDDeleteNotFound{}
}

// WithPayload adds the payload to the delete origin location Id delete not found response
func (o *DeleteOriginLocationIDDeleteNotFound) WithPayload(payload *models.Error) *DeleteOriginLocationIDDeleteNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete origin location Id delete not found response
func (o *DeleteOriginLocationIDDeleteNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteOriginLocationIDDeleteNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteOriginLocationIDDeleteDefault Unexpected error

swagger:response deleteOriginLocationIdDeleteDefault
*/
type DeleteOriginLocationIDDeleteDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteOriginLocationIDDeleteDefault creates DeleteOriginLocationIDDeleteDefault with default headers values
func NewDeleteOriginLocationIDDeleteDefault(code int) *DeleteOriginLocationIDDeleteDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteOriginLocationIDDeleteDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete origin location ID delete default response
func (o *DeleteOriginLocationIDDeleteDefault) WithStatusCode(code int) *DeleteOriginLocationIDDeleteDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete origin location ID delete default response
func (o *DeleteOriginLocationIDDeleteDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete origin location ID delete default response
func (o *DeleteOriginLocationIDDeleteDefault) WithPayload(payload *models.Error) *DeleteOriginLocationIDDeleteDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete origin location ID delete default response
func (o *DeleteOriginLocationIDDeleteDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteOriginLocationIDDeleteDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
