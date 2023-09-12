// Code generated by go-swagger; DO NOT EDIT.

package invites

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostInviteIDAcceptNoContentCode is the HTTP code returned for type PostInviteIDAcceptNoContent
const PostInviteIDAcceptNoContentCode int = 204

/*PostInviteIDAcceptNoContent the invite was accepted

swagger:response postInviteIdAcceptNoContent
*/
type PostInviteIDAcceptNoContent struct {
}

// NewPostInviteIDAcceptNoContent creates PostInviteIDAcceptNoContent with default headers values
func NewPostInviteIDAcceptNoContent() *PostInviteIDAcceptNoContent {

	return &PostInviteIDAcceptNoContent{}
}

// WriteResponse to the client
func (o *PostInviteIDAcceptNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostInviteIDAcceptForbiddenCode is the HTTP code returned for type PostInviteIDAcceptForbidden
const PostInviteIDAcceptForbiddenCode int = 403

/*PostInviteIDAcceptForbidden access forbidden

swagger:response postInviteIdAcceptForbidden
*/
type PostInviteIDAcceptForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostInviteIDAcceptForbidden creates PostInviteIDAcceptForbidden with default headers values
func NewPostInviteIDAcceptForbidden() *PostInviteIDAcceptForbidden {

	return &PostInviteIDAcceptForbidden{}
}

// WithPayload adds the payload to the post invite Id accept forbidden response
func (o *PostInviteIDAcceptForbidden) WithPayload(payload *models.Error) *PostInviteIDAcceptForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post invite Id accept forbidden response
func (o *PostInviteIDAcceptForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostInviteIDAcceptForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostInviteIDAcceptNotFoundCode is the HTTP code returned for type PostInviteIDAcceptNotFound
const PostInviteIDAcceptNotFoundCode int = 404

/*PostInviteIDAcceptNotFound cant find invite

swagger:response postInviteIdAcceptNotFound
*/
type PostInviteIDAcceptNotFound struct {
}

// NewPostInviteIDAcceptNotFound creates PostInviteIDAcceptNotFound with default headers values
func NewPostInviteIDAcceptNotFound() *PostInviteIDAcceptNotFound {

	return &PostInviteIDAcceptNotFound{}
}

// WriteResponse to the client
func (o *PostInviteIDAcceptNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// PostInviteIDAcceptConflictCode is the HTTP code returned for type PostInviteIDAcceptConflict
const PostInviteIDAcceptConflictCode int = 409

/*PostInviteIDAcceptConflict cant accept invite

swagger:response postInviteIdAcceptConflict
*/
type PostInviteIDAcceptConflict struct {
}

// NewPostInviteIDAcceptConflict creates PostInviteIDAcceptConflict with default headers values
func NewPostInviteIDAcceptConflict() *PostInviteIDAcceptConflict {

	return &PostInviteIDAcceptConflict{}
}

// WriteResponse to the client
func (o *PostInviteIDAcceptConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(409)
}

/*PostInviteIDAcceptDefault Unexpected error

swagger:response postInviteIdAcceptDefault
*/
type PostInviteIDAcceptDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostInviteIDAcceptDefault creates PostInviteIDAcceptDefault with default headers values
func NewPostInviteIDAcceptDefault(code int) *PostInviteIDAcceptDefault {
	if code <= 0 {
		code = 500
	}

	return &PostInviteIDAcceptDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post invite ID accept default response
func (o *PostInviteIDAcceptDefault) WithStatusCode(code int) *PostInviteIDAcceptDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post invite ID accept default response
func (o *PostInviteIDAcceptDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post invite ID accept default response
func (o *PostInviteIDAcceptDefault) WithPayload(payload *models.Error) *PostInviteIDAcceptDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post invite ID accept default response
func (o *PostInviteIDAcceptDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostInviteIDAcceptDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}