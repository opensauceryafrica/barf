package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {

	// barf tries to be as unobtrusive as possible, so your route handlers still
	// inherit the standard http.ResponseWriter and *http.Request parameters
	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Hello World",
		})
	})

	// create & start server
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}
}
