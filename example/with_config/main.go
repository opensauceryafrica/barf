package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {
	// create server
	allow := true
	if err := barf.Stark(barf.Augment{
		Port:     "5000",
		Logging:  &allow, // enable request logging
		Recovery: &allow, // enable panic recovery so barf returns a 500 error instead of crashing
	}); err != nil {
		log.Fatal(err)
	}

	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Hello World",
		})
	})

	// start barf server
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}
}
