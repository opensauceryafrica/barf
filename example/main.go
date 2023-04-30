package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {

	type Env struct {
		// Port for the server to listen on
		Port string `barfenv:"key=PORT;required=true"`
	}

	env := new(Env) // global environment variable

	// load environment variables
	if err := barf.Env(env, "example/.env"); err != nil {
		log.Fatal(err)
	}

	// create server
	if err := barf.Stark(barf.Augment{
		Port:    env.Port,
		Logging: true, // enable request logging
	}); err != nil {
		log.Fatal(err)
	}

	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Status(w, http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Hello World",
		})
	})

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}

}
