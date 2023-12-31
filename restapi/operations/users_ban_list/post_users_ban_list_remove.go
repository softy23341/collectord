// Code generated by go-swagger; DO NOT EDIT.

package users_ban_list

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostUsersBanListRemoveHandlerFunc turns a function with the right signature into a post users ban list remove handler
type PostUsersBanListRemoveHandlerFunc func(PostUsersBanListRemoveParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostUsersBanListRemoveHandlerFunc) Handle(params PostUsersBanListRemoveParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostUsersBanListRemoveHandler interface for that can handle valid post users ban list remove params
type PostUsersBanListRemoveHandler interface {
	Handle(PostUsersBanListRemoveParams, interface{}) middleware.Responder
}

// NewPostUsersBanListRemove creates a new http.Handler for the post users ban list remove operation
func NewPostUsersBanListRemove(ctx *middleware.Context, handler PostUsersBanListRemoveHandler) *PostUsersBanListRemove {
	return &PostUsersBanListRemove{Context: ctx, Handler: handler}
}

/*PostUsersBanListRemove swagger:route POST /users/ban-list/remove Users ban list postUsersBanListRemove

Remove user from block list

*/
type PostUsersBanListRemove struct {
	Context *middleware.Context
	Handler PostUsersBanListRemoveHandler
}

func (o *PostUsersBanListRemove) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostUsersBanListRemoveParams()

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
