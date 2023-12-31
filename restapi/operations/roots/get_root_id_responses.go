// Code generated by go-swagger; DO NOT EDIT.

package roots

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// GetRootIDOKCode is the HTTP code returned for type GetRootIDOK
const GetRootIDOKCode int = 200

/*GetRootIDOK success resonse

swagger:response getRootIdOK
*/
type GetRootIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.ARoot `json:"body,omitempty"`
}

// NewGetRootIDOK creates GetRootIDOK with default headers values
func NewGetRootIDOK() *GetRootIDOK {

	return &GetRootIDOK{}
}

// WithPayload adds the payload to the get root Id o k response
func (o *GetRootIDOK) WithPayload(payload *models.ARoot) *GetRootIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get root Id o k response
func (o *GetRootIDOK) SetPayload(payload *models.ARoot) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRootIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRootIDForbiddenCode is the HTTP code returned for type GetRootIDForbidden
const GetRootIDForbiddenCode int = 403

/*GetRootIDForbidden Forbidden

swagger:response getRootIdForbidden
*/
type GetRootIDForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRootIDForbidden creates GetRootIDForbidden with default headers values
func NewGetRootIDForbidden() *GetRootIDForbidden {

	return &GetRootIDForbidden{}
}

// WithPayload adds the payload to the get root Id forbidden response
func (o *GetRootIDForbidden) WithPayload(payload *models.Error) *GetRootIDForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get root Id forbidden response
func (o *GetRootIDForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRootIDForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetRootIDDefault Unexpected error

swagger:response getRootIdDefault
*/
type GetRootIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRootIDDefault creates GetRootIDDefault with default headers values
func NewGetRootIDDefault(code int) *GetRootIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetRootIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get root ID default response
func (o *GetRootIDDefault) WithStatusCode(code int) *GetRootIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get root ID default response
func (o *GetRootIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get root ID default response
func (o *GetRootIDDefault) WithPayload(payload *models.Error) *GetRootIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get root ID default response
func (o *GetRootIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRootIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
