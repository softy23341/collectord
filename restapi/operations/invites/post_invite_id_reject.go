// Code generated by go-swagger; DO NOT EDIT.

package invites

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostInviteIDRejectHandlerFunc turns a function with the right signature into a post invite ID reject handler
type PostInviteIDRejectHandlerFunc func(PostInviteIDRejectParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostInviteIDRejectHandlerFunc) Handle(params PostInviteIDRejectParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostInviteIDRejectHandler interface for that can handle valid post invite ID reject params
type PostInviteIDRejectHandler interface {
	Handle(PostInviteIDRejectParams, interface{}) middleware.Responder
}

// NewPostInviteIDReject creates a new http.Handler for the post invite ID reject operation
func NewPostInviteIDReject(ctx *middleware.Context, handler PostInviteIDRejectHandler) *PostInviteIDReject {
	return &PostInviteIDReject{Context: ctx, Handler: handler}
}

/*PostInviteIDReject swagger:route POST /invite/{ID}/reject Invites postInviteIdReject

PostInviteIDReject post invite ID reject API

*/
type PostInviteIDReject struct {
	Context *middleware.Context
	Handler PostInviteIDRejectHandler
}

func (o *PostInviteIDReject) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostInviteIDRejectParams()

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