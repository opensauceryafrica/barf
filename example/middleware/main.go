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
				return origin == "https://www.google.com"
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	// apply middleware
	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("before 0")
			h.ServeHTTP(w, r)
			barf.Logger().Info("after 0")
		})
	})

	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("before 1")
			h.ServeHTTP(w, r)
			barf.Logger().Info("after 1")
		})
	})

	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("before 2")
			h.ServeHTTP(w, r)
			barf.Logger().Info("after 2")
		})
	})

	barf.Get("/dashboard/:username", func(w http.ResponseWriter, r *http.Request) {
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
