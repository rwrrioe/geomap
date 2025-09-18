package main

import (
	"log"

	server "github.com/rwrrioe/geomap/backend"
)

func main() {
	server := server.NewHTTPServer()

	err := server.InitServerDefault()
	if err != nil {
		log.Fatal(err)
	}
}
