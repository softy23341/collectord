// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// GetUserAboutURL generates an URL for the get user about operation
type GetUserAboutURL struct {
	ID     *int64
	RootID *int64

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetUserAboutURL) WithBasePath(bp string) *GetUserAboutURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetUserAboutURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetUserAboutURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/user/about"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/v1"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var id string
	if o.ID != nil {
		id = swag.FormatInt64(*o.ID)
	}
	if id != "" {
		qs.Set("ID", id)
	}

	var rootID string
	if o.RootID != nil {
		rootID = swag.FormatInt64(*o.RootID)
	}
	if rootID != "" {
		qs.Set("rootID", rootID)
	}

	result.RawQuery = qs.Encode()

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetUserAboutURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetUserAboutURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetUserAboutURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetUserAboutURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetUserAboutURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetUserAboutURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
