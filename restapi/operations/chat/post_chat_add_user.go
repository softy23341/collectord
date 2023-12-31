// Code generated by go-swagger; DO NOT EDIT.

package chat

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostChatAddUserHandlerFunc turns a function with the right signature into a post chat add user handler
type PostChatAddUserHandlerFunc func(PostChatAddUserParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostChatAddUserHandlerFunc) Handle(params PostChatAddUserParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostChatAddUserHandler interface for that can handle valid post chat add user params
type PostChatAddUserHandler interface {
	Handle(PostChatAddUserParams, interface{}) middleware.Responder
}

// NewPostChatAddUser creates a new http.Handler for the post chat add user operation
func NewPostChatAddUser(ctx *middleware.Context, handler PostChatAddUserHandler) *PostChatAddUser {
	return &PostChatAddUser{Context: ctx, Handler: handler}
}

/*PostChatAddUser swagger:route POST /chat/add-user Chat postChatAddUser

Add user to chat

*/
type PostChatAddUser struct {
	Context *middleware.Context
	Handler PostChatAddUserHandler
}

func (o *PostChatAddUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostChatAddUserParams()

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
