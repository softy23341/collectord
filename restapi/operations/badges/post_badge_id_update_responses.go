// Code generated by go-swagger; DO NOT EDIT.

package badges

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostBadgeIDUpdateNoContentCode is the HTTP code returned for type PostBadgeIDUpdateNoContent
const PostBadgeIDUpdateNoContentCode int = 204

/*PostBadgeIDUpdateNoContent success

swagger:response postBadgeIdUpdateNoContent
*/
type PostBadgeIDUpdateNoContent struct {
}

// NewPostBadgeIDUpdateNoContent creates PostBadgeIDUpdateNoContent with default headers values
func NewPostBadgeIDUpdateNoContent() *PostBadgeIDUpdateNoContent {

	return &PostBadgeIDUpdateNoContent{}
}

// WriteResponse to the client
func (o *PostBadgeIDUpdateNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostBadgeIDUpdateForbiddenCode is the HTTP code returned for type PostBadgeIDUpdateForbidden
const PostBadgeIDUpdateForbiddenCode int = 403

/*PostBadgeIDUpdateForbidden Forbidden

swagger:response postBadgeIdUpdateForbidden
*/
type PostBadgeIDUpdateForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostBadgeIDUpdateForbidden creates PostBadgeIDUpdateForbidden with default headers values
func NewPostBadgeIDUpdateForbidden() *PostBadgeIDUpdateForbidden {

	return &PostBadgeIDUpdateForbidden{}
}

// WithPayload adds the payload to the post badge Id update forbidden response
func (o *PostBadgeIDUpdateForbidden) WithPayload(payload *models.Error) *PostBadgeIDUpdateForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post badge Id update forbidden response
func (o *PostBadgeIDUpdateForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBadgeIDUpdateForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostBadgeIDUpdateNotFoundCode is the HTTP code returned for type PostBadgeIDUpdateNotFound
const PostBadgeIDUpdateNotFoundCode int = 404

/*PostBadgeIDUpdateNotFound cant find badge

swagger:response postBadgeIdUpdateNotFound
*/
type PostBadgeIDUpdateNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostBadgeIDUpdateNotFound creates PostBadgeIDUpdateNotFound with default headers values
func NewPostBadgeIDUpdateNotFound() *PostBadgeIDUpdateNotFound {

	return &PostBadgeIDUpdateNotFound{}
}

// WithPayload adds the payload to the post badge Id update not found response
func (o *PostBadgeIDUpdateNotFound) WithPayload(payload *models.Error) *PostBadgeIDUpdateNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post badge Id update not found response
func (o *PostBadgeIDUpdateNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBadgeIDUpdateNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostBadgeIDUpdateConflictCode is the HTTP code returned for type PostBadgeIDUpdateConflict
const PostBadgeIDUpdateConflictCode int = 409

/*PostBadgeIDUpdateConflict Badge already present (color and name must be uniq)

swagger:response postBadgeIdUpdateConflict
*/
type PostBadgeIDUpdateConflict struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostBadgeIDUpdateConflict creates PostBadgeIDUpdateConflict with default headers values
func NewPostBadgeIDUpdateConflict() *PostBadgeIDUpdateConflict {

	return &PostBadgeIDUpdateConflict{}
}

// WithPayload adds the payload to the post badge Id update conflict response
func (o *PostBadgeIDUpdateConflict) WithPayload(payload *models.Error) *PostBadgeIDUpdateConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post badge Id update conflict response
func (o *PostBadgeIDUpdateConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBadgeIDUpdateConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostBadgeIDUpdateDefault Unexpected error

swagger:response postBadgeIdUpdateDefault
*/
type PostBadgeIDUpdateDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostBadgeIDUpdateDefault creates PostBadgeIDUpdateDefault with default headers values
func NewPostBadgeIDUpdateDefault(code int) *PostBadgeIDUpdateDefault {
	if code <= 0 {
		code = 500
	}

	return &PostBadgeIDUpdateDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post badge ID update default response
func (o *PostBadgeIDUpdateDefault) WithStatusCode(code int) *PostBadgeIDUpdateDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post badge ID update default response
func (o *PostBadgeIDUpdateDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post badge ID update default response
func (o *PostBadgeIDUpdateDefault) WithPayload(payload *models.Error) *PostBadgeIDUpdateDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post badge ID update default response
func (o *PostBadgeIDUpdateDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBadgeIDUpdateDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
