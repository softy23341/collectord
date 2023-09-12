// Code generated by go-swagger; DO NOT EDIT.

package collections

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostCollectionsObjectsHandlerFunc turns a function with the right signature into a post collections objects handler
type PostCollectionsObjectsHandlerFunc func(PostCollectionsObjectsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCollectionsObjectsHandlerFunc) Handle(params PostCollectionsObjectsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PostCollectionsObjectsHandler interface for that can handle valid post collections objects params
type PostCollectionsObjectsHandler interface {
	Handle(PostCollectionsObjectsParams, interface{}) middleware.Responder
}

// NewPostCollectionsObjects creates a new http.Handler for the post collections objects operation
func NewPostCollectionsObjects(ctx *middleware.Context, handler PostCollectionsObjectsHandler) *PostCollectionsObjects {
	return &PostCollectionsObjects{Context: ctx, Handler: handler}
}

/*PostCollectionsObjects swagger:route POST /collections/objects Collections postCollectionsObjects

Collection objects

*/
type PostCollectionsObjects struct {
	Context *middleware.Context
	Handler PostCollectionsObjectsHandler
}

func (o *PostCollectionsObjects) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostCollectionsObjectsParams()

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