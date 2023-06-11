package main

import (
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
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// apply middleware to all routes
	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("before 1")
			h.ServeHTTP(w, r)
			barf.Logger().Info("after 1")
		})
	})

	// create a subrouter (retroframe)
	var r *barf.SubRoute = barf.RetroFrame("/api")

	// apply middleware to subrouter. note that the only difference between this and global middlewares is that you need to pass the
	barf.Hippocampus(r).Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("sub before 0")
			h.ServeHTTP(w, r)
			barf.Logger().Info("sub after 0")
		})
	})

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Home",
		})
	})

	// create another subrouter (retroframe)
	// note that although the path is the same, the subroute is different and won't inherit the middleware from the previous subroute
	s := barf.RetroFrame("/api").RetroFrame("/v1")
	s.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "About",
		})
	})

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("This is a middleware for the contact route")
			next.ServeHTTP(w, r)
			barf.Logger().Info("This is a middleware for the contact route after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("This is the second middleware for the contact route")
			next.ServeHTTP(w, r)
			barf.Logger().Info("This is the second middleware for the contact route after")
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			barf.Logger().Info("This is the third middleware for the contact route")
			next.ServeHTTP(w, r)
			barf.Logger().Info("This is the third middleware for the contact route after")
		})
	}	

	// create another subrouter from the previous subrouter
	// note that, in this case, the subroute will inherit the middleware from the previous subroute
	n := r.RetroFrame("/v2")
	n.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Contact",
		})
	}).Use(middleware, middleware2)

	x := r.RetroFrame("/v3")
	x.Use(middleware3)
	x.Get("/contacts", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Contact version 3",
		})
	})

	x.Get("/contactz", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Contactz version 3",
		})
	})

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
