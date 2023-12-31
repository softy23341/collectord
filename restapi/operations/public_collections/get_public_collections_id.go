// Code generated by go-swagger; DO NOT EDIT.

package public_collections

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetPublicCollectionsIDHandlerFunc turns a function with the right signature into a get public collections ID handler
type GetPublicCollectionsIDHandlerFunc func(GetPublicCollectionsIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPublicCollectionsIDHandlerFunc) Handle(params GetPublicCollectionsIDParams) middleware.Responder {
	return fn(params)
}

// GetPublicCollectionsIDHandler interface for that can handle valid get public collections ID params
type GetPublicCollectionsIDHandler interface {
	Handle(GetPublicCollectionsIDParams) middleware.Responder
}

// NewGetPublicCollectionsID creates a new http.Handler for the get public collections ID operation
func NewGetPublicCollectionsID(ctx *middleware.Context, handler GetPublicCollectionsIDHandler) *GetPublicCollectionsID {
	return &GetPublicCollectionsID{Context: ctx, Handler: handler}
}

/*GetPublicCollectionsID swagger:route GET /public/collections/{ID} Public Collections getPublicCollectionsId

Public collection info

*/
type GetPublicCollectionsID struct {
	Context *middleware.Context
	Handler GetPublicCollectionsIDHandler
}

func (o *GetPublicCollectionsID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPublicCollectionsIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
