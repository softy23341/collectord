// Code generated by go-swagger; DO NOT EDIT.

package messages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostMessageHandlerFunc turns a function with the right signature into a post message handler
type PostMessageHandlerFunc func(PostMessageParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostMessageHandlerFunc) Handle(params PostMessageParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostMessageHandler interface for that can handle valid post message params
type PostMessageHandler interface {
	Handle(PostMessageParams, interface{}) middleware.Responder
}

// NewPostMessage creates a new http.Handler for the post message operation
func NewPostMessage(ctx *middleware.Context, handler PostMessageHandler) *PostMessage {
	return &PostMessage{Context: ctx, Handler: handler}
}

/*PostMessage swagger:route POST /message Messages postMessage

Send message

*/
type PostMessage struct {
	Context *middleware.Context
	Handler PostMessageHandler
}

func (o *PostMessage) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostMessageParams()

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