// Code generated by go-swagger; DO NOT EDIT.

package originlocations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteOriginLocationIDDeleteHandlerFunc turns a function with the right signature into a delete origin location ID delete handler
type DeleteOriginLocationIDDeleteHandlerFunc func(DeleteOriginLocationIDDeleteParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteOriginLocationIDDeleteHandlerFunc) Handle(params DeleteOriginLocationIDDeleteParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteOriginLocationIDDeleteHandler interface for that can handle valid delete origin location ID delete params
type DeleteOriginLocationIDDeleteHandler interface {
	Handle(DeleteOriginLocationIDDeleteParams, interface{}) middleware.Responder
}

// NewDeleteOriginLocationIDDelete creates a new http.Handler for the delete origin location ID delete operation
func NewDeleteOriginLocationIDDelete(ctx *middleware.Context, handler DeleteOriginLocationIDDeleteHandler) *DeleteOriginLocationIDDelete {
	return &DeleteOriginLocationIDDelete{Context: ctx, Handler: handler}
}

/*DeleteOriginLocationIDDelete swagger:route DELETE /origin-location/{ID}/delete Originlocations deleteOriginLocationIdDelete

DeleteOriginLocationIDDelete delete origin location ID delete API

*/
type DeleteOriginLocationIDDelete struct {
	Context *middleware.Context
	Handler DeleteOriginLocationIDDeleteHandler
}

func (o *DeleteOriginLocationIDDelete) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteOriginLocationIDDeleteParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
