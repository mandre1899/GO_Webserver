package main

import (
	"github.com/mandre1899/GO_Webserver/internal/server"
)

func main() {
	server := server.WebServer{}

	server.CreateServer()

	server.Server.ListenAndServe()
}
