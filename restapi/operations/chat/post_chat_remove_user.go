// Code generated by go-swagger; DO NOT EDIT.

package chat

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostChatRemoveUserHandlerFunc turns a function with the right signature into a post chat remove user handler
type PostChatRemoveUserHandlerFunc func(PostChatRemoveUserParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostChatRemoveUserHandlerFunc) Handle(params PostChatRemoveUserParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostChatRemoveUserHandler interface for that can handle valid post chat remove user params
type PostChatRemoveUserHandler interface {
	Handle(PostChatRemoveUserParams, interface{}) middleware.Responder
}

// NewPostChatRemoveUser creates a new http.Handler for the post chat remove user operation
func NewPostChatRemoveUser(ctx *middleware.Context, handler PostChatRemoveUserHandler) *PostChatRemoveUser {
	return &PostChatRemoveUser{Context: ctx, Handler: handler}
}

/*PostChatRemoveUser swagger:route POST /chat/remove-user Chat postChatRemoveUser

Remove user from chat

*/
type PostChatRemoveUser struct {
	Context *middleware.Context
	Handler PostChatRemoveUserHandler
}

func (o *PostChatRemoveUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostChatRemoveUserParams()

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
