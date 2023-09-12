// Code generated by go-swagger; DO NOT EDIT.

package session

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostSessionRegisterDeviceTokenHandlerFunc turns a function with the right signature into a post session register device token handler
type PostSessionRegisterDeviceTokenHandlerFunc func(PostSessionRegisterDeviceTokenParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostSessionRegisterDeviceTokenHandlerFunc) Handle(params PostSessionRegisterDeviceTokenParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostSessionRegisterDeviceTokenHandler interface for that can handle valid post session register device token params
type PostSessionRegisterDeviceTokenHandler interface {
	Handle(PostSessionRegisterDeviceTokenParams, interface{}) middleware.Responder
}

// NewPostSessionRegisterDeviceToken creates a new http.Handler for the post session register device token operation
func NewPostSessionRegisterDeviceToken(ctx *middleware.Context, handler PostSessionRegisterDeviceTokenHandler) *PostSessionRegisterDeviceToken {
	return &PostSessionRegisterDeviceToken{Context: ctx, Handler: handler}
}

/*PostSessionRegisterDeviceToken swagger:route POST /session/register-device-token Session postSessionRegisterDeviceToken

Register device

*/
type PostSessionRegisterDeviceToken struct {
	Context *middleware.Context
	Handler PostSessionRegisterDeviceTokenHandler
}

func (o *PostSessionRegisterDeviceToken) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostSessionRegisterDeviceTokenParams()

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