package main

import (
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func main() {
	// store all the storm connections here.
	allStorms := []*barf.Storm{}

	// barf tries to be as unobtrusive as possible, so your route handlers still
	// inherit the standard http.ResponseWriter and *http.Request parameters
	barf.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		barf.Logger().Info("reached endpoint")
		stormBreaker := barf.StormBreaker{}
		storm, err := stormBreaker.Upgrade(w, r, nil)
		if err != nil {
			barf.Logger().Error(err.Error())
			return
		}

		// add every new storm connection into the array to hold all clients in the websocket pool
		allStorms = append(allStorms, storm)
		barf.Logger().Info("upgraded to websocket")
		for {
			barf.Logger().Info("checking for websocket messages...")
			messageType, message, err := storm.ReadMessage()
			if err != nil {
				barf.Logger().Error(err.Error())
				return
			}

			// write the received message to all storms in the websocket connection.
			for _, stormUsers := range allStorms {
				err = stormUsers.WriteMessage(messageType, message)
				if err != nil {
					barf.Logger().Error(err.Error())
					return
				}
			}
		}
	})

	// create & start server
	if err := barf.Beck(); err != nil {
		// barf exposes a logger instance
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
