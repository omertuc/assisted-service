// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DownloadClusterISOHeadersCloneHandlerFunc turns a function with the right signature into a download cluster i s o headers clone handler
type DownloadClusterISOHeadersCloneHandlerFunc func(DownloadClusterISOHeadersCloneParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DownloadClusterISOHeadersCloneHandlerFunc) Handle(params DownloadClusterISOHeadersCloneParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DownloadClusterISOHeadersCloneHandler interface for that can handle valid download cluster i s o headers clone params
type DownloadClusterISOHeadersCloneHandler interface {
	Handle(DownloadClusterISOHeadersCloneParams, interface{}) middleware.Responder
}

// NewDownloadClusterISOHeadersClone creates a new http.Handler for the download cluster i s o headers clone operation
func NewDownloadClusterISOHeadersClone(ctx *middleware.Context, handler DownloadClusterISOHeadersCloneHandler) *DownloadClusterISOHeadersClone {
	return &DownloadClusterISOHeadersClone{Context: ctx, Handler: handler}
}

/*DownloadClusterISOHeadersClone swagger:route HEAD /clusters/{cluster_id}/downloads/image.iso installer downloadClusterISOHeadersClone

Downloads the OpenShift per-cluster Discovery ISO Headers only.

*/
type DownloadClusterISOHeadersClone struct {
	Context *middleware.Context
	Handler DownloadClusterISOHeadersCloneHandler
}

func (o *DownloadClusterISOHeadersClone) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDownloadClusterISOHeadersCloneParams()

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
