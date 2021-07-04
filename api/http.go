package api

import (
	"net/url"
	"path"
)

// queryParams holds URL query parameters in a structured way
type queryParams map[string]string

func (a *API) constructWsURL(p string) string {
	return a.constructURL("ws", p, nil)
}

func (a *API) constructHttpURL(p string, q queryParams) string {
	return a.constructURL(a.baseURL.Scheme, p, q)
}

func (a *API) constructURL(s, p string, q queryParams) string {
	u := url.URL{Scheme: s, Host: a.baseURL.Host, Path: path.Join("/"+Version, p)}
	// set query parameters while preserving any previous ones
	if q != nil {
		original := u.Query()
		for p, v := range q {
			original.Set(p, v)
		}
		u.RawQuery = original.Encode()
	}
	return u.String()
}
