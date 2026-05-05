package server

import (
	"fmt"
	"net/http"
	"database/sql"
	"os"
	"github.com/mandre1899/GO_Webserver/internal/database"
	"github.com/mandre1899/GO_Webserver/internal/api"
	"github.com/mandre1899/GO_Webserver/internal/middleware"
)

type WebServer struct {
	Server	*http.Server
	Mux		*http.ServeMux
	ApiConf	middleware.ApiConfig
	Db		*database.Queries
}

func (s *WebServer) CreateServer() {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("Connection to db failed %v", err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
	    fmt.Print("Database ping failed %v", err)
	    os.Exit(1)
	}
	s.Db = database.New(db)

	mux := http.NewServeMux()
	s.ApiConf = middleware.ApiConfig{}
	s.Server = &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	
	api.RegisterRoutes(mux, s.Db, &s.ApiConf)
}

