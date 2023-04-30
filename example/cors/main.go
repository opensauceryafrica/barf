package main

import (
	"log"
	"net/http"
	"time"

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
	logging := true
	if err := barf.Stark(barf.Augment{
		Port:    env.Port,
		Logging: &logging, // enable request logging
		CORS: &barf.CORS{
			AllowedOrigins: []string{"https://*.google.com"},
			MaxAge:         3600,
			AllowedMethods: []string{
				http.MethodGet,
			},
			AllowedOriginFunc: func(origin string) bool {
				if origin == "secure.com" {
					return false
				} else if origin == "insecure.com" {
					return true
				}
				return false
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	barf.Get("/dashboard/:username", func(w http.ResponseWriter, r *http.Request) {
		<-time.After(2 * time.Second)
		params, _ := barf.Request(r).Params().JSON()
		query, _ := barf.Request(r).Query().JSON()
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    map[string]interface{}{"params": params, "query": query},
			Message: "Hello World",
		})
	})

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}
}
