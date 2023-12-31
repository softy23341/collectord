// Code generated by go-swagger; DO NOT EDIT.

package medias

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// GetMediasNoContentCode is the HTTP code returned for type GetMediasNoContent
const GetMediasNoContentCode int = 204

/*GetMediasNoContent success with X-Accel-Redirect

swagger:response getMediasNoContent
*/
type GetMediasNoContent struct {
}

// NewGetMediasNoContent creates GetMediasNoContent with default headers values
func NewGetMediasNoContent() *GetMediasNoContent {

	return &GetMediasNoContent{}
}

// WriteResponse to the client
func (o *GetMediasNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// GetMediasForbiddenCode is the HTTP code returned for type GetMediasForbidden
const GetMediasForbiddenCode int = 403

/*GetMediasForbidden access denied

swagger:response getMediasForbidden
*/
type GetMediasForbidden struct {
}

// NewGetMediasForbidden creates GetMediasForbidden with default headers values
func NewGetMediasForbidden() *GetMediasForbidden {

	return &GetMediasForbidden{}
}

// WriteResponse to the client
func (o *GetMediasForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

/*GetMediasDefault Unexpected error

swagger:response getMediasDefault
*/
type GetMediasDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetMediasDefault creates GetMediasDefault with default headers values
func NewGetMediasDefault(code int) *GetMediasDefault {
	if code <= 0 {
		code = 500
	}

	return &GetMediasDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get medias default response
func (o *GetMediasDefault) WithStatusCode(code int) *GetMediasDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get medias default response
func (o *GetMediasDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get medias default response
func (o *GetMediasDefault) WithPayload(payload *models.Error) *GetMediasDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get medias default response
func (o *GetMediasDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMediasDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
