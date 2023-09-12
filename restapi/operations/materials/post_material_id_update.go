// Code generated by go-swagger; DO NOT EDIT.

package materials

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostMaterialIDUpdateHandlerFunc turns a function with the right signature into a post material ID update handler
type PostMaterialIDUpdateHandlerFunc func(PostMaterialIDUpdateParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostMaterialIDUpdateHandlerFunc) Handle(params PostMaterialIDUpdateParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostMaterialIDUpdateHandler interface for that can handle valid post material ID update params
type PostMaterialIDUpdateHandler interface {
	Handle(PostMaterialIDUpdateParams, interface{}) middleware.Responder
}

// NewPostMaterialIDUpdate creates a new http.Handler for the post material ID update operation
func NewPostMaterialIDUpdate(ctx *middleware.Context, handler PostMaterialIDUpdateHandler) *PostMaterialIDUpdate {
	return &PostMaterialIDUpdate{Context: ctx, Handler: handler}
}

/*PostMaterialIDUpdate swagger:route POST /material/{ID}/update Materials postMaterialIdUpdate

update material

*/
type PostMaterialIDUpdate struct {
	Context *middleware.Context
	Handler PostMaterialIDUpdateHandler
}

func (o *PostMaterialIDUpdate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostMaterialIDUpdateParams()

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
