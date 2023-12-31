// Code generated by go-swagger; DO NOT EDIT.

package badges

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetBadgeGetHandlerFunc turns a function with the right signature into a get badge get handler
type GetBadgeGetHandlerFunc func(GetBadgeGetParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBadgeGetHandlerFunc) Handle(params GetBadgeGetParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetBadgeGetHandler interface for that can handle valid get badge get params
type GetBadgeGetHandler interface {
	Handle(GetBadgeGetParams, interface{}) middleware.Responder
}

// NewGetBadgeGet creates a new http.Handler for the get badge get operation
func NewGetBadgeGet(ctx *middleware.Context, handler GetBadgeGetHandler) *GetBadgeGet {
	return &GetBadgeGet{Context: ctx, Handler: handler}
}

/*GetBadgeGet swagger:route GET /badge/get Badges getBadgeGet

Get badges for root

*/
type GetBadgeGet struct {
	Context *middleware.Context
	Handler GetBadgeGetHandler
}

func (o *GetBadgeGet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetBadgeGetParams()

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
