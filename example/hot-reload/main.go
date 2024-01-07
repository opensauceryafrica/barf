package main

import (
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func main() {
	// create server
	if err := barf.Stark(barf.Augment{
		Logging:        barf.Allow(), // enable request logging
		Recovery:       barf.Allow(), // enable panic recovery so barf returns a 500 error instead of crashing
		AllowHotReload: barf.Allow(), // allow hot reload
		HotReload: &barf.HotReload{
			Root:        ".",
			ExcludeDir:  []string{},
			IncludeDir:  []string{},
			StopOnError: false,
			BuildCmd:    "go build -o app",
			Bin:         "./app",
		},
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
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
		// barf exposes a logger instance
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
