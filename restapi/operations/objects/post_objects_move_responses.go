// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostObjectsMoveNoContentCode is the HTTP code returned for type PostObjectsMoveNoContent
const PostObjectsMoveNoContentCode int = 204

/*PostObjectsMoveNoContent Success

swagger:response postObjectsMoveNoContent
*/
type PostObjectsMoveNoContent struct {
}

// NewPostObjectsMoveNoContent creates PostObjectsMoveNoContent with default headers values
func NewPostObjectsMoveNoContent() *PostObjectsMoveNoContent {

	return &PostObjectsMoveNoContent{}
}

// WriteResponse to the client
func (o *PostObjectsMoveNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostObjectsMoveForbiddenCode is the HTTP code returned for type PostObjectsMoveForbidden
const PostObjectsMoveForbiddenCode int = 403

/*PostObjectsMoveForbidden Forbidden

swagger:response postObjectsMoveForbidden
*/
type PostObjectsMoveForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostObjectsMoveForbidden creates PostObjectsMoveForbidden with default headers values
func NewPostObjectsMoveForbidden() *PostObjectsMoveForbidden {

	return &PostObjectsMoveForbidden{}
}

// WithPayload adds the payload to the post objects move forbidden response
func (o *PostObjectsMoveForbidden) WithPayload(payload *models.Error) *PostObjectsMoveForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post objects move forbidden response
func (o *PostObjectsMoveForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostObjectsMoveForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostObjectsMoveNotFoundCode is the HTTP code returned for type PostObjectsMoveNotFound
const PostObjectsMoveNotFoundCode int = 404

/*PostObjectsMoveNotFound cant find objects

swagger:response postObjectsMoveNotFound
*/
type PostObjectsMoveNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostObjectsMoveNotFound creates PostObjectsMoveNotFound with default headers values
func NewPostObjectsMoveNotFound() *PostObjectsMoveNotFound {

	return &PostObjectsMoveNotFound{}
}

// WithPayload adds the payload to the post objects move not found response
func (o *PostObjectsMoveNotFound) WithPayload(payload *models.Error) *PostObjectsMoveNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post objects move not found response
func (o *PostObjectsMoveNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostObjectsMoveNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostObjectsMoveDefault Unexpected error

swagger:response postObjectsMoveDefault
*/
type PostObjectsMoveDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostObjectsMoveDefault creates PostObjectsMoveDefault with default headers values
func NewPostObjectsMoveDefault(code int) *PostObjectsMoveDefault {
	if code <= 0 {
		code = 500
	}

	return &PostObjectsMoveDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post objects move default response
func (o *PostObjectsMoveDefault) WithStatusCode(code int) *PostObjectsMoveDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post objects move default response
func (o *PostObjectsMoveDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post objects move default response
func (o *PostObjectsMoveDefault) WithPayload(payload *models.Error) *PostObjectsMoveDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post objects move default response
func (o *PostObjectsMoveDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostObjectsMoveDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
