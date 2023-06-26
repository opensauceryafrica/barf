/*
	package barf

Basically, A Remarkable Framework!
*/
package barf

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/opensaucerer/barf/constant"
	logger "github.com/opensaucerer/barf/log"
	"github.com/opensaucerer/barf/router"
	"github.com/opensaucerer/barf/server"
	"github.com/opensaucerer/barf/typing"
)

func createServer(a typing.Augment) error {

	// create handler
	server.Mux = http.NewServeMux()

	var r http.Handler = server.Mux

	// wrap router into logger middleware
	if *server.Augment.Logging {
		r = logger.Morgan(r)
	}

	// wrap router into router middleware
	r = router.Router(server.JSON)(r)

	// wrap router into cors middleware
	// r = middleware.CORS(middleware.Prepare(*server.Augment.CORS))(r)

	// // wrap router into recover middleware
	// if *server.Augment.Recovery {
	// 	r = middleware.Recover(server.JSON)(r)
	// }

	// create barf for hijacking
	router.Barf.Router = r
	router.Barf.Stack = []typing.Middleware{}

	// create server
	server.HTTP = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", server.Augment.Host, server.Augment.Port),
		ReadTimeout:       time.Duration(server.Augment.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(server.Augment.WriteTimeout) * time.Second,
		MaxHeaderBytes:    server.Augment.MaxHeaderBytes,
		Handler:           r,
		ReadHeaderTimeout: time.Duration(server.Augment.ReadHeaderTimeout) * time.Second,
	}

	// this will load the CORS and Recovery middleware into the stack
	if *server.Augment.Recovery {
		Hippocampus().Hijack()
		logger.Info("Recovery middleware added to base barf handler")
	}

	return nil
}

// Stark retrieves any existing barf server or creates a new one and returns an error, if any.
// You can optionally pass in a barf.Augment struct to override the default config.
// To start the server, call the bart.Beck()
func Stark(augmentation ...typing.Augment) error {
	// return nil if server already exists
	if server.HTTP != nil {
		return nil
	}
	augu := typing.Augment{
		MaxHeaderBytes:    constant.MaxHeaderBytes,
		ReadTimeout:       constant.ReadTimeout,
		ReadHeaderTimeout: constant.ReadTimeout,
		WriteTimeout:      constant.WriteTimeout,
		ShutdownTimeout:   constant.ShutdownTimeout,
		Host:              constant.Host,
		Port:              constant.Port,
		UseHTTPS:          constant.UseHTTPS,
		Logging:           &constant.Logging,
		Recovery:          &constant.Recovery,
		CORS:              &typing.CORS{},
	}
	if len(augmentation) > 0 {
		// validate the struct
		t := reflect.TypeOf(augmentation[0])
		if t.Kind() != reflect.Struct {
			return fmt.Errorf("error: Stark() expects a struct, got %s", t.Kind())
		}
		// validate struct is a barf.Augment
		if t.Name() != "Augment" {
			return fmt.Errorf("error: Stark() expects a barf.Augment struct, got %s", t.Name())
		}
		// override the default config
		aug := augmentation[0]
		// load default configurations
		if aug.MaxHeaderBytes != 0 {
			augu.MaxHeaderBytes = aug.MaxHeaderBytes
		}
		if aug.ReadTimeout != 0 {
			augu.ReadTimeout = aug.ReadTimeout
		}
		if aug.WriteTimeout != 0 {
			augu.WriteTimeout = aug.WriteTimeout
		}
		if aug.Host != "" {
			augu.Host = aug.Host
		}
		if aug.Port != "" {
			augu.Port = aug.Port
		}
		if aug.ReadHeaderTimeout != 0 {
			augu.ReadHeaderTimeout = aug.ReadHeaderTimeout
		}
		if aug.Logging != nil {
			augu.Logging = aug.Logging
		}
		if aug.Recovery != nil {
			augu.Recovery = aug.Recovery
		}
		if aug.CORS != nil {
			augu.CORS = aug.CORS
		}
		if aug.UseHTTPS {
			if aug.SSLCertFile == "" || aug.SSLKeyFile == "" {
				return fmt.Errorf("error: Stark() enabling UseHTTPS requires SSLCertFile & SSLKeyFile")
			}
			augu.SSLCertFile = aug.SSLCertFile
			augu.SSLKeyFile = aug.SSLKeyFile
		}
		augu.UseHTTPS = aug.UseHTTPS
	}
	// make config global
	server.Augment = &augu
	return createServer(augu)
}

// Beck starts the barf server and returns an error, if any. Alternatively, Beck also creates a new barf server with the default config and starts it, only if barf.Stark() was not called before.
func Beck() error {
	// return nil if server already Beckoned
	if server.Beckoned != nil && *server.Beckoned {
		return nil
	}
	// if barf.Stark() was not called, call it
	if server.HTTP == nil {
		if err := Stark(); err != nil {
			return err
		}
	}
	// register shutdown function
	go func() {
		valid := true
		server.Beckoned = &valid
		shutdown()
	}()

	// start server with https enabled
	if server.Augment.UseHTTPS {
		logger.Info(fmt.Sprintf("BARF server started at https://%s:%s", server.Augment.Host, server.Augment.Port))
		if err := server.HTTP.ListenAndServeTLS(server.Augment.SSLCertFile, server.Augment.SSLKeyFile); err != nil {
			server.Beckoned = nil
			return err
		}
	} else {
		// start server
		logger.Info(fmt.Sprintf("BARF server started at http://%s:%s", server.Augment.Host, server.Augment.Port))
		if err := server.HTTP.ListenAndServe(); err != nil {
			server.Beckoned = nil
			return err
		}
	}
	return nil
}

// shutdown gracefully shuts down the server with the specified timeout.
func shutdown() {
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need to add it
	signal.Notify(constant.ShutdownChan, syscall.SIGINT, syscall.SIGTERM)
	<-constant.ShutdownChan
	logger.Warn("Shutting down BARF...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(server.Augment.ShutdownTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.HTTP.Shutdown(ctx); err != nil {
		logger.Error("BARF forced to shut down...")
		log.Fatal()
	}
	logger.Debug("BARF exited!")
}
