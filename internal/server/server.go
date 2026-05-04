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
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=uft-8")
		w.Write([]byte("OK"))
	})
}

