// Code generated by go-swagger; DO NOT EDIT.

package public_collections

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// PostPublicCollectionsObjectsOKCode is the HTTP code returned for type PostPublicCollectionsObjectsOK
const PostPublicCollectionsObjectsOKCode int = 200

/*PostPublicCollectionsObjectsOK Collection objects list

swagger:response postPublicCollectionsObjectsOK
*/
type PostPublicCollectionsObjectsOK struct {

	/*
	  In: Body
	*/
	Payload *models.AObjectsPreview `json:"body,omitempty"`
}

// NewPostPublicCollectionsObjectsOK creates PostPublicCollectionsObjectsOK with default headers values
func NewPostPublicCollectionsObjectsOK() *PostPublicCollectionsObjectsOK {

	return &PostPublicCollectionsObjectsOK{}
}

// WithPayload adds the payload to the post public collections objects o k response
func (o *PostPublicCollectionsObjectsOK) WithPayload(payload *models.AObjectsPreview) *PostPublicCollectionsObjectsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post public collections objects o k response
func (o *PostPublicCollectionsObjectsOK) SetPayload(payload *models.AObjectsPreview) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostPublicCollectionsObjectsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostPublicCollectionsObjectsForbiddenCode is the HTTP code returned for type PostPublicCollectionsObjectsForbidden
const PostPublicCollectionsObjectsForbiddenCode int = 403

/*PostPublicCollectionsObjectsForbidden Forbidden

swagger:response postPublicCollectionsObjectsForbidden
*/
type PostPublicCollectionsObjectsForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostPublicCollectionsObjectsForbidden creates PostPublicCollectionsObjectsForbidden with default headers values
func NewPostPublicCollectionsObjectsForbidden() *PostPublicCollectionsObjectsForbidden {

	return &PostPublicCollectionsObjectsForbidden{}
}

// WithPayload adds the payload to the post public collections objects forbidden response
func (o *PostPublicCollectionsObjectsForbidden) WithPayload(payload *models.Error) *PostPublicCollectionsObjectsForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post public collections objects forbidden response
func (o *PostPublicCollectionsObjectsForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostPublicCollectionsObjectsForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostPublicCollectionsObjectsDefault Unexpected error

swagger:response postPublicCollectionsObjectsDefault
*/
type PostPublicCollectionsObjectsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostPublicCollectionsObjectsDefault creates PostPublicCollectionsObjectsDefault with default headers values
func NewPostPublicCollectionsObjectsDefault(code int) *PostPublicCollectionsObjectsDefault {
	if code <= 0 {
		code = 500
	}

	return &PostPublicCollectionsObjectsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post public collections objects default response
func (o *PostPublicCollectionsObjectsDefault) WithStatusCode(code int) *PostPublicCollectionsObjectsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post public collections objects default response
func (o *PostPublicCollectionsObjectsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post public collections objects default response
func (o *PostPublicCollectionsObjectsDefault) WithPayload(payload *models.Error) *PostPublicCollectionsObjectsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post public collections objects default response
func (o *PostPublicCollectionsObjectsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostPublicCollectionsObjectsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}