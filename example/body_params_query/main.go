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
	allow := true
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  &allow,
		Recovery: &allow,
	}); err != nil {
		log.Fatal(err)
	}

	barf.Post("/:username", func(w http.ResponseWriter, r *http.Request) {

		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := barf.Request(r).Body().Format(&data)
		if err != nil {
			barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
				Status:  false,
				Data:    nil,
				Message: "Invalid request body",
			})
			return
		}

		params, _ := barf.Request(r).Params().JSON()
		query, _ := barf.Request(r).Query().JSON()

		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    map[string]interface{}{"params": params, "query": query, "body": data},
			Message: "Hello World",
		})
	})

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}
}
