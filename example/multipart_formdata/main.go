package main

import (
	"io"
	"net/http"
	"os"

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
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// create server
	allow := true
	disallow := false
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  &allow, // enable request logging
		Recovery: &disallow,
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
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// create a subrouter (retroframe)
	s := barf.RetroFrame("/api").RetroFrame("/v1")
	s.Get("/about", func(w http.ResponseWriter, r *http.Request) {

		message := "About"

		// parsing form-data
		body, err := barf.Request(r).Form().Body().JSON()
		if err != nil {
			message = err.Error()
		}

		head := barf.Request(r).Form().File().Get("file")
		file, _ := head.Open()
		defer file.Close()

		// save file
		f, err := os.OpenFile(head.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			message = err.Error()
		}
		defer f.Close()
		io.Copy(f, file)

		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    body,
			Message: message,
		})
	})

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
