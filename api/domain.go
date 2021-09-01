package api

import "github.com/cloudfauj/cloudfauj/server"

func (a *API) AddDomain(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/domain/"+name+"/add"), nil)
}

func (a *API) DeleteDomain(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/domain/"+name+"/delete"), nil)
}
