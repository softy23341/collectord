// Code generated by go-swagger; DO NOT EDIT.

package roots

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostRootAddUserNoContentCode is the HTTP code returned for type PostRootAddUserNoContent
const PostRootAddUserNoContentCode int = 204

/*PostRootAddUserNoContent success

swagger:response postRootAddUserNoContent
*/
type PostRootAddUserNoContent struct {
}

// NewPostRootAddUserNoContent creates PostRootAddUserNoContent with default headers values
func NewPostRootAddUserNoContent() *PostRootAddUserNoContent {

	return &PostRootAddUserNoContent{}
}

// WriteResponse to the client
func (o *PostRootAddUserNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// PostRootAddUserForbiddenCode is the HTTP code returned for type PostRootAddUserForbidden
const PostRootAddUserForbiddenCode int = 403

/*PostRootAddUserForbidden Forbidden

swagger:response postRootAddUserForbidden
*/
type PostRootAddUserForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRootAddUserForbidden creates PostRootAddUserForbidden with default headers values
func NewPostRootAddUserForbidden() *PostRootAddUserForbidden {

	return &PostRootAddUserForbidden{}
}

// WithPayload adds the payload to the post root add user forbidden response
func (o *PostRootAddUserForbidden) WithPayload(payload *models.Error) *PostRootAddUserForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post root add user forbidden response
func (o *PostRootAddUserForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRootAddUserForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostRootAddUserDefault Unexpected error

swagger:response postRootAddUserDefault
*/
type PostRootAddUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostRootAddUserDefault creates PostRootAddUserDefault with default headers values
func NewPostRootAddUserDefault(code int) *PostRootAddUserDefault {
	if code <= 0 {
		code = 500
	}

	return &PostRootAddUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post root add user default response
func (o *PostRootAddUserDefault) WithStatusCode(code int) *PostRootAddUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post root add user default response
func (o *PostRootAddUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post root add user default response
func (o *PostRootAddUserDefault) WithPayload(payload *models.Error) *PostRootAddUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post root add user default response
func (o *PostRootAddUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostRootAddUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
