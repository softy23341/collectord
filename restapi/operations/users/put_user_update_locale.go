// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PutUserUpdateLocaleHandlerFunc turns a function with the right signature into a put user update locale handler
type PutUserUpdateLocaleHandlerFunc func(PutUserUpdateLocaleParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutUserUpdateLocaleHandlerFunc) Handle(params PutUserUpdateLocaleParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutUserUpdateLocaleHandler interface for that can handle valid put user update locale params
type PutUserUpdateLocaleHandler interface {
	Handle(PutUserUpdateLocaleParams, interface{}) middleware.Responder
}

// NewPutUserUpdateLocale creates a new http.Handler for the put user update locale operation
func NewPutUserUpdateLocale(ctx *middleware.Context, handler PutUserUpdateLocaleHandler) *PutUserUpdateLocale {
	return &PutUserUpdateLocale{Context: ctx, Handler: handler}
}

/*PutUserUpdateLocale swagger:route PUT /user/update/locale Users putUserUpdateLocale

edit user locale

*/
type PutUserUpdateLocale struct {
	Context *middleware.Context
	Handler PutUserUpdateLocaleHandler
}

func (o *PutUserUpdateLocale) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutUserUpdateLocaleParams()

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
