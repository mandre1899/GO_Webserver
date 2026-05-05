package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mandre1899/GO_Webserver/internal/database"
	"github.com/mandre1899/GO_Webserver/internal/middleware"
)

func RegisterRoutes(mux *http.ServeMux, db *database.Queries, apiConf *middleware.ApiConfig) {
	mux.Handle("/app/", apiConf.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		hits := apiConf.FileserverHits.Load()
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
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiConf.FileserverHits.Store(0)
		db.DeleteAllUsers(context.Background())
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=uft-8")
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/chirps", ValidateChirpHandler(db, apiConf))
	mux.HandleFunc("GET /api/chirps", GetChrips(db))
	mux.HandleFunc("GET /api/chirps/{id}", GetChripById(db))
	mux.HandleFunc("POST /api/users", CreateUser(db))
	mux.HandleFunc("POST /api/login", LoginUser(db, apiConf))
}

