// Code generated by go-swagger; DO NOT EDIT.

package materials

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "git.softndit.com/collector/backend/models"
)

// DeleteMaterialIDDeleteNoContentCode is the HTTP code returned for type DeleteMaterialIDDeleteNoContent
const DeleteMaterialIDDeleteNoContentCode int = 204

/*DeleteMaterialIDDeleteNoContent succes

swagger:response deleteMaterialIdDeleteNoContent
*/
type DeleteMaterialIDDeleteNoContent struct {
}

// NewDeleteMaterialIDDeleteNoContent creates DeleteMaterialIDDeleteNoContent with default headers values
func NewDeleteMaterialIDDeleteNoContent() *DeleteMaterialIDDeleteNoContent {

	return &DeleteMaterialIDDeleteNoContent{}
}

// WriteResponse to the client
func (o *DeleteMaterialIDDeleteNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteMaterialIDDeleteForbiddenCode is the HTTP code returned for type DeleteMaterialIDDeleteForbidden
const DeleteMaterialIDDeleteForbiddenCode int = 403

/*DeleteMaterialIDDeleteForbidden access forbidden

swagger:response deleteMaterialIdDeleteForbidden
*/
type DeleteMaterialIDDeleteForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteMaterialIDDeleteForbidden creates DeleteMaterialIDDeleteForbidden with default headers values
func NewDeleteMaterialIDDeleteForbidden() *DeleteMaterialIDDeleteForbidden {

	return &DeleteMaterialIDDeleteForbidden{}
}

// WithPayload adds the payload to the delete material Id delete forbidden response
func (o *DeleteMaterialIDDeleteForbidden) WithPayload(payload *models.Error) *DeleteMaterialIDDeleteForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete material Id delete forbidden response
func (o *DeleteMaterialIDDeleteForbidden) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteMaterialIDDeleteForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteMaterialIDDeleteNotFoundCode is the HTTP code returned for type DeleteMaterialIDDeleteNotFound
const DeleteMaterialIDDeleteNotFoundCode int = 404

/*DeleteMaterialIDDeleteNotFound cant find material

swagger:response deleteMaterialIdDeleteNotFound
*/
type DeleteMaterialIDDeleteNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteMaterialIDDeleteNotFound creates DeleteMaterialIDDeleteNotFound with default headers values
func NewDeleteMaterialIDDeleteNotFound() *DeleteMaterialIDDeleteNotFound {

	return &DeleteMaterialIDDeleteNotFound{}
}

// WithPayload adds the payload to the delete material Id delete not found response
func (o *DeleteMaterialIDDeleteNotFound) WithPayload(payload *models.Error) *DeleteMaterialIDDeleteNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete material Id delete not found response
func (o *DeleteMaterialIDDeleteNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteMaterialIDDeleteNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteMaterialIDDeleteDefault Unexpected error

swagger:response deleteMaterialIdDeleteDefault
*/
type DeleteMaterialIDDeleteDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteMaterialIDDeleteDefault creates DeleteMaterialIDDeleteDefault with default headers values
func NewDeleteMaterialIDDeleteDefault(code int) *DeleteMaterialIDDeleteDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteMaterialIDDeleteDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete material ID delete default response
func (o *DeleteMaterialIDDeleteDefault) WithStatusCode(code int) *DeleteMaterialIDDeleteDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete material ID delete default response
func (o *DeleteMaterialIDDeleteDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete material ID delete default response
func (o *DeleteMaterialIDDeleteDefault) WithPayload(payload *models.Error) *DeleteMaterialIDDeleteDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete material ID delete default response
func (o *DeleteMaterialIDDeleteDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteMaterialIDDeleteDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}