# barf - Basically A Remarkable Framework

BARF is a small and unobtrusive framework for building JSON-based web APIs on REST or GraphQL-based architectures. It is implemented such that getting started is easy and quick, but it is also flexible enough to allow for more complex use cases.

- No application instance
- No bullsh\*t context
- ...and no re-inventing the wheel

It’s just what you need being provided to you in an unobtrusive way.

![image](https://user-images.githubusercontent.com/59074379/236047064-a16887b2-31f6-4e66-bf6b-340fa1b18a53.png)

## Installation

```shell
go get github.com/opensaucerer/barf
```

## Usage

For a comprehensive overview on how to use barf, please refer to the [example](./example) folder and its sub folders.

### A simple BARF REST API

```go
package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {

	// barf tries to be as unobtrusive as possible, so your route handlers still
	// inherit the standard http.ResponseWriter and *http.Request parameters
	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Hello World",
		})
	})

	// create & start server
	if err := barf.Beck(); err != nil {
		// barf exposes a logger instance
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
```

### BARF with custom configuration

```go
package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {
	// create server
	if err := barf.Stark(barf.Augment{
		Port:     "5000",
		Logging:  barf.Allow(),  // enable request logging
		Recovery: barf.Allow(), // enable panic recovery so barf returns a 500 error instead of crashing
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
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
```

### BARF with Environment Variables

```go
package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {

	type Env struct {
		// Port for the server to listen on
		Port string `barfenv:"key=PORT;required=true"` // barfenv tag allows barf to load environment variables
	}

	env := new(Env) // global environment variable

	// you can use barf to dynamically load environment variables into a struct
	if err := barf.Env(env, "example/.env"); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// create server
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  barf.Allow(),
		Recovery: barf.Allow(),
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

	// start server - create & start server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
```

### BARF with request body, variable paths and query parameters

```go
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
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// create server
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  barf.Allow(),
		Recovery: barf.Allow(),
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
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
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}

```

### BARF with multipart/form-data

```go
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
	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  barf.Allow(), // enable request logging
		Recovery: barf.Disallow(),
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
```

### BARF with custom middleware

```go
package main

import (
	"log"
	"net/http"

	"github.com/opensaucerer/barf"
)

func main() {
	// create server
	if err := barf.Stark(barf.Augment{
		Port:     "5000",
		Logging:  barf.Allow(),  // enable request logging
		Recovery: barf.Allow(), // enable panic recovery so barf returns a 500 error instead of crashing
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	// apply global middleware to all routes - middleware is applied in the order it is added and must be added before the call to barf.Beck()
	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("before 0")
			h.ServeHTTP(w, r)
			log.Println("after 0")
		})
	})

	// define routes
	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Hello World",
		})
	})

	// apply another global middleware
	barf.Hippocampus().Hijack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("before 1")
			h.ServeHTTP(w, r)
			log.Println("after 1")
		})
	})

	// start barf server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
```

### BARF with subrouters

```go
package main

import (
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func main() {
	// create server
	if err := barf.Stark(barf.Augment{
		Port:     "5000",
		Logging:  barf.Allow(), // enable request logging
		Recovery: barf.Allow(), // enable panic recovery so barf returns a 500 error instead of crashing
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

	// create a subrouter (retroframe)
	s := barf.RetroFrame("/api/v1")
	s.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "About",
		})
	})

	// start barf server
	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
```

# Contributing

Barf is an open source project and we welcome contributions of all kinds to help improve the project. Please read our [contributing guide](./contributing.md) to learn about our development process, how to propose bug fixes and improvements, and how to build and test your changes to Barf.

For a starting point on features Barf currently lacks, see the [issues page](https://github.com/opensaucerer/barf/issues).

# License

Barf is [MIT licensed](./LICENSE).

# Contributors

<a href="https://github.com/opensaucerer/barf/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=opensaucerer/barf" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
