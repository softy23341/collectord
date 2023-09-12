// Code generated by go-swagger; DO NOT EDIT.

package roots

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostRootRemoveUserNoContentCode is the HTTP code returned for type PostRootRemoveUserNoContent
const PostRootRemoveUserNoContentCode int = 204

/*PostRootRemoveUserNoContent success

swagger:response postRootRemoveUserNoContent
*/
type PostRootRemoveUserNoContent struct {
}

// NewPostRootRemoveUserNoContent creates PostRootRemoveUserNoContent with default headers values
func NewPostRootRemoveUserNoContent() *PostRootRemoveUserNoContent {

	return &PostRootRemoveUserNoContent{}
}

// WriteResponse to the client
func (o *PostRootRemoveUserNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostRootRemoveUserForbiddenCode is the HTTP code returned for type PostRootRemoveUserForbidden
const PostRootRemoveUserForbiddenCode int = 403

/*PostRootRemoveUserForbidden Forbidden

swagger:response postRootRemoveUserForbidden
*/
type PostRootRemoveUserForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRootRemoveUserForbidden creates PostRootRemoveUserForbidden with default headers values
func NewPostRootRemoveUserForbidden() *PostRootRemoveUserForbidden {

	return &PostRootRemoveUserForbidden{}
}

// WithPayload adds the payload to the post root remove user forbidden response
func (o *PostRootRemoveUserForbidden) WithPayload(payload *models.Error) *PostRootRemoveUserForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post root remove user forbidden response
func (o *PostRootRemoveUserForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRootRemoveUserForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostRootRemoveUserNotFoundCode is the HTTP code returned for type PostRootRemoveUserNotFound
const PostRootRemoveUserNotFoundCode int = 404

/*PostRootRemoveUserNotFound cant find user or root

swagger:response postRootRemoveUserNotFound
*/
type PostRootRemoveUserNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRootRemoveUserNotFound creates PostRootRemoveUserNotFound with default headers values
func NewPostRootRemoveUserNotFound() *PostRootRemoveUserNotFound {

	return &PostRootRemoveUserNotFound{}
}

// WithPayload adds the payload to the post root remove user not found response
func (o *PostRootRemoveUserNotFound) WithPayload(payload *models.Error) *PostRootRemoveUserNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post root remove user not found response
func (o *PostRootRemoveUserNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRootRemoveUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostRootRemoveUserDefault Unexpected error

swagger:response postRootRemoveUserDefault
*/
type PostRootRemoveUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRootRemoveUserDefault creates PostRootRemoveUserDefault with default headers values
func NewPostRootRemoveUserDefault(code int) *PostRootRemoveUserDefault {
	if code <= 0 {
		code = 500
	}

	return &PostRootRemoveUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post root remove user default response
func (o *PostRootRemoveUserDefault) WithStatusCode(code int) *PostRootRemoveUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post root remove user default response
func (o *PostRootRemoveUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post root remove user default response
func (o *PostRootRemoveUserDefault) WithPayload(payload *models.Error) *PostRootRemoveUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post root remove user default response
func (o *PostRootRemoveUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRootRemoveUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}