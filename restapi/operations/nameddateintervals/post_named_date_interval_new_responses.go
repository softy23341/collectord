// Code generated by go-swagger; DO NOT EDIT.

package nameddateintervals

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostNamedDateIntervalNewOKCode is the HTTP code returned for type PostNamedDateIntervalNewOK
const PostNamedDateIntervalNewOKCode int = 200

/*PostNamedDateIntervalNewOK Success

swagger:response postNamedDateIntervalNewOK
*/
type PostNamedDateIntervalNewOK struct {

	/*
	  In: Body
	*/
	Payload *models.ACreateNamedDateInterval `json:"body,omitempty"`
}

// NewPostNamedDateIntervalNewOK creates PostNamedDateIntervalNewOK with default headers values
func NewPostNamedDateIntervalNewOK() *PostNamedDateIntervalNewOK {

	return &PostNamedDateIntervalNewOK{}
}

// WithPayload adds the payload to the post named date interval new o k response
func (o *PostNamedDateIntervalNewOK) WithPayload(payload *models.ACreateNamedDateInterval) *PostNamedDateIntervalNewOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post named date interval new o k response
func (o *PostNamedDateIntervalNewOK) SetPayload(payload *models.ACreateNamedDateInterval) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostNamedDateIntervalNewOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostNamedDateIntervalNewForbiddenCode is the HTTP code returned for type PostNamedDateIntervalNewForbidden
const PostNamedDateIntervalNewForbiddenCode int = 403

/*PostNamedDateIntervalNewForbidden access forbidden

swagger:response postNamedDateIntervalNewForbidden
*/
type PostNamedDateIntervalNewForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostNamedDateIntervalNewForbidden creates PostNamedDateIntervalNewForbidden with default headers values
func NewPostNamedDateIntervalNewForbidden() *PostNamedDateIntervalNewForbidden {

	return &PostNamedDateIntervalNewForbidden{}
}

// WithPayload adds the payload to the post named date interval new forbidden response
func (o *PostNamedDateIntervalNewForbidden) WithPayload(payload *models.Error) *PostNamedDateIntervalNewForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post named date interval new forbidden response
func (o *PostNamedDateIntervalNewForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostNamedDateIntervalNewForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostNamedDateIntervalNewConflictCode is the HTTP code returned for type PostNamedDateIntervalNewConflict
const PostNamedDateIntervalNewConflictCode int = 409

/*PostNamedDateIntervalNewConflict named date interval already present

swagger:response postNamedDateIntervalNewConflict
*/
type PostNamedDateIntervalNewConflict struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostNamedDateIntervalNewConflict creates PostNamedDateIntervalNewConflict with default headers values
func NewPostNamedDateIntervalNewConflict() *PostNamedDateIntervalNewConflict {

	return &PostNamedDateIntervalNewConflict{}
}

// WithPayload adds the payload to the post named date interval new conflict response
func (o *PostNamedDateIntervalNewConflict) WithPayload(payload *models.Error) *PostNamedDateIntervalNewConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post named date interval new conflict response
func (o *PostNamedDateIntervalNewConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostNamedDateIntervalNewConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostNamedDateIntervalNewDefault Unexpected error

swagger:response postNamedDateIntervalNewDefault
*/
type PostNamedDateIntervalNewDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostNamedDateIntervalNewDefault creates PostNamedDateIntervalNewDefault with default headers values
func NewPostNamedDateIntervalNewDefault(code int) *PostNamedDateIntervalNewDefault {
	if code <= 0 {
		code = 500
	}

	return &PostNamedDateIntervalNewDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post named date interval new default response
func (o *PostNamedDateIntervalNewDefault) WithStatusCode(code int) *PostNamedDateIntervalNewDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post named date interval new default response
func (o *PostNamedDateIntervalNewDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post named date interval new default response
func (o *PostNamedDateIntervalNewDefault) WithPayload(payload *models.Error) *PostNamedDateIntervalNewDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post named date interval new default response
func (o *PostNamedDateIntervalNewDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostNamedDateIntervalNewDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}