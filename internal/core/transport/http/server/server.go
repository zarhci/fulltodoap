package core_server

import "net/http"

type Server struct {
	mux *http.ServeMux
}
