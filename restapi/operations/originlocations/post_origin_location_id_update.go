// Code generated by go-swagger; DO NOT EDIT.

package originlocations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostOriginLocationIDUpdateHandlerFunc turns a function with the right signature into a post origin location ID update handler
type PostOriginLocationIDUpdateHandlerFunc func(PostOriginLocationIDUpdateParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostOriginLocationIDUpdateHandlerFunc) Handle(params PostOriginLocationIDUpdateParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostOriginLocationIDUpdateHandler interface for that can handle valid post origin location ID update params
type PostOriginLocationIDUpdateHandler interface {
	Handle(PostOriginLocationIDUpdateParams, interface{}) middleware.Responder
}

// NewPostOriginLocationIDUpdate creates a new http.Handler for the post origin location ID update operation
func NewPostOriginLocationIDUpdate(ctx *middleware.Context, handler PostOriginLocationIDUpdateHandler) *PostOriginLocationIDUpdate {
	return &PostOriginLocationIDUpdate{Context: ctx, Handler: handler}
}

/*PostOriginLocationIDUpdate swagger:route POST /origin-location/{ID}/update Originlocations postOriginLocationIdUpdate

update origin-location

*/
type PostOriginLocationIDUpdate struct {
	Context *middleware.Context
	Handler PostOriginLocationIDUpdateHandler
}

func (o *PostOriginLocationIDUpdate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostOriginLocationIDUpdateParams()

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
