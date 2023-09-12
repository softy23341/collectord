// Code generated by go-swagger; DO NOT EDIT.

package nameddateintervals

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteNamedDateIntervalIDDeleteHandlerFunc turns a function with the right signature into a delete named date interval ID delete handler
type DeleteNamedDateIntervalIDDeleteHandlerFunc func(DeleteNamedDateIntervalIDDeleteParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteNamedDateIntervalIDDeleteHandlerFunc) Handle(params DeleteNamedDateIntervalIDDeleteParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteNamedDateIntervalIDDeleteHandler interface for that can handle valid delete named date interval ID delete params
type DeleteNamedDateIntervalIDDeleteHandler interface {
	Handle(DeleteNamedDateIntervalIDDeleteParams, interface{}) middleware.Responder
}

// NewDeleteNamedDateIntervalIDDelete creates a new http.Handler for the delete named date interval ID delete operation
func NewDeleteNamedDateIntervalIDDelete(ctx *middleware.Context, handler DeleteNamedDateIntervalIDDeleteHandler) *DeleteNamedDateIntervalIDDelete {
	return &DeleteNamedDateIntervalIDDelete{Context: ctx, Handler: handler}
}

/*DeleteNamedDateIntervalIDDelete swagger:route DELETE /named-date-interval/{ID}/delete Nameddateintervals deleteNamedDateIntervalIdDelete

DeleteNamedDateIntervalIDDelete delete named date interval ID delete API

*/
type DeleteNamedDateIntervalIDDelete struct {
	Context *middleware.Context
	Handler DeleteNamedDateIntervalIDDeleteHandler
}

func (o *DeleteNamedDateIntervalIDDelete) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteNamedDateIntervalIDDeleteParams()

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