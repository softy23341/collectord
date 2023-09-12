// Code generated by go-swagger; DO NOT EDIT.

package invites

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostInviteNewHandlerFunc turns a function with the right signature into a post invite new handler
type PostInviteNewHandlerFunc func(PostInviteNewParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostInviteNewHandlerFunc) Handle(params PostInviteNewParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostInviteNewHandler interface for that can handle valid post invite new params
type PostInviteNewHandler interface {
	Handle(PostInviteNewParams, interface{}) middleware.Responder
}

// NewPostInviteNew creates a new http.Handler for the post invite new operation
func NewPostInviteNew(ctx *middleware.Context, handler PostInviteNewHandler) *PostInviteNew {
	return &PostInviteNew{Context: ctx, Handler: handler}
}

/*PostInviteNew swagger:route POST /invite/new Invites postInviteNew

PostInviteNew post invite new API

*/
type PostInviteNew struct {
	Context *middleware.Context
	Handler PostInviteNewHandler
}

func (o *PostInviteNew) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostInviteNewParams()

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
