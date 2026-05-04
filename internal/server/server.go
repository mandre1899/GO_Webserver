package server

import (
	"fmt"
	"net/http"

	"github.com/mandre1899/GO_Webserver/internal/middleware"
	"github.com/mandre1899/GO_Webserver/internal/api"
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
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		hits := s.ApiConf.FileserverHits.Load()
		resStr := fmt.Sprintf(`
			<html>
			  <body>
			    <h1>Welcome, Chirpy Admin</h1>
			    <p>Chirpy has been visited %d times!</p>
			  </body>
			</html>
		`, hits)
		w.Write([]byte(resStr))
	})
	mux.HandleFunc("POST /admin/reset", func (w http.ResponseWriter, r *http.Request)  {
		s.ApiConf.FileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=uft-8")
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirpHandler)
}

