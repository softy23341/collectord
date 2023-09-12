// Code generated by go-swagger; DO NOT EDIT.

package collections

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetCollectionsIDHandlerFunc turns a function with the right signature into a get collections ID handler
type GetCollectionsIDHandlerFunc func(GetCollectionsIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCollectionsIDHandlerFunc) Handle(params GetCollectionsIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetCollectionsIDHandler interface for that can handle valid get collections ID params
type GetCollectionsIDHandler interface {
	Handle(GetCollectionsIDParams, interface{}) middleware.Responder
}

// NewGetCollectionsID creates a new http.Handler for the get collections ID operation
func NewGetCollectionsID(ctx *middleware.Context, handler GetCollectionsIDHandler) *GetCollectionsID {
	return &GetCollectionsID{Context: ctx, Handler: handler}
}

/*GetCollectionsID swagger:route GET /collections/{ID} Collections getCollectionsId

Collection info

*/
type GetCollectionsID struct {
	Context *middleware.Context
	Handler GetCollectionsIDHandler
}

func (o *GetCollectionsID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetCollectionsIDParams()

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