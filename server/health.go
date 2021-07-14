package server

import "net/http"

func (s *server) handlerGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
