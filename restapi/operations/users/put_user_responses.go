// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PutUserNoContentCode is the HTTP code returned for type PutUserNoContent
const PutUserNoContentCode int = 204

/*PutUserNoContent success

swagger:response putUserNoContent
*/
type PutUserNoContent struct {
}

// NewPutUserNoContent creates PutUserNoContent with default headers values
func NewPutUserNoContent() *PutUserNoContent {

	return &PutUserNoContent{}
}

// WriteResponse to the client
func (o *PutUserNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PutUserForbiddenCode is the HTTP code returned for type PutUserForbidden
const PutUserForbiddenCode int = 403

/*PutUserForbidden Forbidden

swagger:response putUserForbidden
*/
type PutUserForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutUserForbidden creates PutUserForbidden with default headers values
func NewPutUserForbidden() *PutUserForbidden {

	return &PutUserForbidden{}
}

// WithPayload adds the payload to the put user forbidden response
func (o *PutUserForbidden) WithPayload(payload *models.Error) *PutUserForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put user forbidden response
func (o *PutUserForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutUserForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PutUserDefault Unexpected error

swagger:response putUserDefault
*/
type PutUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutUserDefault creates PutUserDefault with default headers values
func NewPutUserDefault(code int) *PutUserDefault {
	if code <= 0 {
		code = 500
	}

	return &PutUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put user default response
func (o *PutUserDefault) WithStatusCode(code int) *PutUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put user default response
func (o *PutUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put user default response
func (o *PutUserDefault) WithPayload(payload *models.Error) *PutUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put user default response
func (o *PutUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
