// Code generated by go-swagger; DO NOT EDIT.

package actors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostActorIDUpdateHandlerFunc turns a function with the right signature into a post actor ID update handler
type PostActorIDUpdateHandlerFunc func(PostActorIDUpdateParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostActorIDUpdateHandlerFunc) Handle(params PostActorIDUpdateParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostActorIDUpdateHandler interface for that can handle valid post actor ID update params
type PostActorIDUpdateHandler interface {
	Handle(PostActorIDUpdateParams, interface{}) middleware.Responder
}

// NewPostActorIDUpdate creates a new http.Handler for the post actor ID update operation
func NewPostActorIDUpdate(ctx *middleware.Context, handler PostActorIDUpdateHandler) *PostActorIDUpdate {
	return &PostActorIDUpdate{Context: ctx, Handler: handler}
}

/*PostActorIDUpdate swagger:route POST /actor/{ID}/update Actors postActorIdUpdate

update actor

*/
type PostActorIDUpdate struct {
	Context *middleware.Context
	Handler PostActorIDUpdateHandler
}

func (o *PostActorIDUpdate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostActorIDUpdateParams()

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
