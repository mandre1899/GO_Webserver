package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mandre1899/GO_Webserver/internal/server"
)

func main() {
	godotenv.Load()
	server := server.WebServer{}

	server.CreateServer()

	server.Server.ListenAndServe()
}
