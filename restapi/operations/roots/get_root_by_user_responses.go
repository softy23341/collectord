// Code generated by go-swagger; DO NOT EDIT.

package roots

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// GetRootByUserOKCode is the HTTP code returned for type GetRootByUserOK
const GetRootByUserOKCode int = 200

/*GetRootByUserOK success resonse

swagger:response getRootByUserOK
*/
type GetRootByUserOK struct {

	/*
	  In: Body
	*/
	Payload *models.ARoots `json:"body,omitempty"`
}

// NewGetRootByUserOK creates GetRootByUserOK with default headers values
func NewGetRootByUserOK() *GetRootByUserOK {

	return &GetRootByUserOK{}
}

// WithPayload adds the payload to the get root by user o k response
func (o *GetRootByUserOK) WithPayload(payload *models.ARoots) *GetRootByUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get root by user o k response
func (o *GetRootByUserOK) SetPayload(payload *models.ARoots) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRootByUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRootByUserForbiddenCode is the HTTP code returned for type GetRootByUserForbidden
const GetRootByUserForbiddenCode int = 403

/*GetRootByUserForbidden Forbidden

swagger:response getRootByUserForbidden
*/
type GetRootByUserForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRootByUserForbidden creates GetRootByUserForbidden with default headers values
func NewGetRootByUserForbidden() *GetRootByUserForbidden {

	return &GetRootByUserForbidden{}
}

// WithPayload adds the payload to the get root by user forbidden response
func (o *GetRootByUserForbidden) WithPayload(payload *models.Error) *GetRootByUserForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get root by user forbidden response
func (o *GetRootByUserForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRootByUserForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetRootByUserDefault Unexpected error

swagger:response getRootByUserDefault
*/
type GetRootByUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRootByUserDefault creates GetRootByUserDefault with default headers values
func NewGetRootByUserDefault(code int) *GetRootByUserDefault {
	if code <= 0 {
		code = 500
	}

	return &GetRootByUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get root by user default response
func (o *GetRootByUserDefault) WithStatusCode(code int) *GetRootByUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get root by user default response
func (o *GetRootByUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get root by user default response
func (o *GetRootByUserDefault) WithPayload(payload *models.Error) *GetRootByUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get root by user default response
func (o *GetRootByUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRootByUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}