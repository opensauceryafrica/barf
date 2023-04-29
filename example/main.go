package main

import (
	"log"
	"regexp"

	"github.com/opensaucerer/barf"
)

func main() {

	path := "/auth/login//"
	log.Println(regexp.MustCompile("^/+|/+$").ReplaceAllString(path, ""))

	// app = Express()

	type Env struct {
		// DatabaseURI is the connection string for database
		DatabaseURI string `barfenv:"key=DATABASE_URI"`
		// Port for the server to listen on
		Port string `barfenv:"key=PORT"`
		// DatabaseName is the database name for the application
		DatabaseName string `barfenv:"key=DATABASE_NAME"`
	}

	env := new(Env) // global environment variable

	// // create server
	// if err := barf.Stark(barf.Augment{
	// 	Port: env.Port,
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	// load environment variables
	if err := barf.Env(env, "example/.env"); err != nil {
		log.Fatal(err)
	}

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		log.Fatal(err)
	}

}
