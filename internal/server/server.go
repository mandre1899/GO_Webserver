package server

import (
	"net/http"
	"strconv"

	"github.com/mandre1899/GO_Webserver/internal/middleware"
)

type WebServer struct {
	Server	*http.Server
	Mux		*http.ServeMux
	ApiConf	middleware.ApiConfig
}

func (s *WebServer) CreateServer() {
	mux := http.NewServeMux()
	s.ApiConf = middleware.ApiConfig{}
	s.Server = &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	mux.Handle("/app/", s.ApiConf.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		hits := s.ApiConf.FileserverHits.Load()
		resStr := "Hits: " + strconv.Itoa(int(hits))
		w.Write([]byte(resStr))
	})
	mux.HandleFunc("POST /reset", func (w http.ResponseWriter, r *http.Request)  {
		s.ApiConf.FileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=uft-8")
		w.Write([]byte("OK"))
	})
}

