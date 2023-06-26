package main

import (
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func main() {

	type Env struct {
		// Port for the server to listen on
		Port       string `barfenv:"key=PORT;required=true"`
		PathToCert string `barfenv:"key=PATH_TO_CERT;required=true"`
		PathToKey  string `barfenv:"key=PATH_TO_KEY;required=true"`
	}

	env := new(Env) // global environment variable

	// load environment variables
	if err := barf.Env(env, "example/.env"); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// create server
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  barf.Allow(),
		Recovery: barf.Allow(),
		// enable TLS for HTTPS
		UseHTTPS:    barf.True(),
		SSLCertFile: env.PathToCert,
		SSLKeyFile:  env.PathToKey,
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// create a subrouter (retroframe)
	var r *barf.SubRoute = barf.RetroFrame("/api")

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {

		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Home",
		})

	})

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
