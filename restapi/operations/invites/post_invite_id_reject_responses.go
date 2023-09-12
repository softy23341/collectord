// Code generated by go-swagger; DO NOT EDIT.

package invites

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostInviteIDRejectNoContentCode is the HTTP code returned for type PostInviteIDRejectNoContent
const PostInviteIDRejectNoContentCode int = 204

/*PostInviteIDRejectNoContent the invite was rejected

swagger:response postInviteIdRejectNoContent
*/
type PostInviteIDRejectNoContent struct {
}

// NewPostInviteIDRejectNoContent creates PostInviteIDRejectNoContent with default headers values
func NewPostInviteIDRejectNoContent() *PostInviteIDRejectNoContent {

	return &PostInviteIDRejectNoContent{}
}

// WriteResponse to the client
func (o *PostInviteIDRejectNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostInviteIDRejectForbiddenCode is the HTTP code returned for type PostInviteIDRejectForbidden
const PostInviteIDRejectForbiddenCode int = 403

/*PostInviteIDRejectForbidden access forbidden

swagger:response postInviteIdRejectForbidden
*/
type PostInviteIDRejectForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostInviteIDRejectForbidden creates PostInviteIDRejectForbidden with default headers values
func NewPostInviteIDRejectForbidden() *PostInviteIDRejectForbidden {

	return &PostInviteIDRejectForbidden{}
}

// WithPayload adds the payload to the post invite Id reject forbidden response
func (o *PostInviteIDRejectForbidden) WithPayload(payload *models.Error) *PostInviteIDRejectForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post invite Id reject forbidden response
func (o *PostInviteIDRejectForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostInviteIDRejectForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostInviteIDRejectNotFoundCode is the HTTP code returned for type PostInviteIDRejectNotFound
const PostInviteIDRejectNotFoundCode int = 404

/*PostInviteIDRejectNotFound cant find invite

swagger:response postInviteIdRejectNotFound
*/
type PostInviteIDRejectNotFound struct {
}

// NewPostInviteIDRejectNotFound creates PostInviteIDRejectNotFound with default headers values
func NewPostInviteIDRejectNotFound() *PostInviteIDRejectNotFound {

	return &PostInviteIDRejectNotFound{}
}

// WriteResponse to the client
func (o *PostInviteIDRejectNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// PostInviteIDRejectConflictCode is the HTTP code returned for type PostInviteIDRejectConflict
const PostInviteIDRejectConflictCode int = 409

/*PostInviteIDRejectConflict cant reject invite

swagger:response postInviteIdRejectConflict
*/
type PostInviteIDRejectConflict struct {
}

// NewPostInviteIDRejectConflict creates PostInviteIDRejectConflict with default headers values
func NewPostInviteIDRejectConflict() *PostInviteIDRejectConflict {

	return &PostInviteIDRejectConflict{}
}

// WriteResponse to the client
func (o *PostInviteIDRejectConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(409)
}

/*PostInviteIDRejectDefault Unexpected error

swagger:response postInviteIdRejectDefault
*/
type PostInviteIDRejectDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostInviteIDRejectDefault creates PostInviteIDRejectDefault with default headers values
func NewPostInviteIDRejectDefault(code int) *PostInviteIDRejectDefault {
	if code <= 0 {
		code = 500
	}

	return &PostInviteIDRejectDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post invite ID reject default response
func (o *PostInviteIDRejectDefault) WithStatusCode(code int) *PostInviteIDRejectDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post invite ID reject default response
func (o *PostInviteIDRejectDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post invite ID reject default response
func (o *PostInviteIDRejectDefault) WithPayload(payload *models.Error) *PostInviteIDRejectDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post invite ID reject default response
func (o *PostInviteIDRejectDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostInviteIDRejectDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
