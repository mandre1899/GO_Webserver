package server

import "net/http"

type WebServer struct {
	Server	*http.Server
	Mux		*http.ServeMux
}

func (s *WebServer) CreateServer() {
	mux := http.NewServeMux()
	s.Server = &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
}

