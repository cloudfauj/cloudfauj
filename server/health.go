package server

import "net/http"

func (s *server) handlerGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	s.log.Info("/health called")
}
