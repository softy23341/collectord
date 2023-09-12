// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostAuthRegConfirmEmailHandlerFunc turns a function with the right signature into a post auth reg confirm email handler
type PostAuthRegConfirmEmailHandlerFunc func(PostAuthRegConfirmEmailParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostAuthRegConfirmEmailHandlerFunc) Handle(params PostAuthRegConfirmEmailParams) middleware.Responder {
	return fn(params)
}

// PostAuthRegConfirmEmailHandler interface for that can handle valid post auth reg confirm email params
type PostAuthRegConfirmEmailHandler interface {
	Handle(PostAuthRegConfirmEmailParams) middleware.Responder
}

// NewPostAuthRegConfirmEmail creates a new http.Handler for the post auth reg confirm email operation
func NewPostAuthRegConfirmEmail(ctx *middleware.Context, handler PostAuthRegConfirmEmailHandler) *PostAuthRegConfirmEmail {
	return &PostAuthRegConfirmEmail{Context: ctx, Handler: handler}
}

/*PostAuthRegConfirmEmail swagger:route POST /auth/reg/confirm-email Auth postAuthRegConfirmEmail

Send token

*/
type PostAuthRegConfirmEmail struct {
	Context *middleware.Context
	Handler PostAuthRegConfirmEmailHandler
}

func (o *PostAuthRegConfirmEmail) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostAuthRegConfirmEmailParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}