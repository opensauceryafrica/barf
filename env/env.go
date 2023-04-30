/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"fmt"
	"reflect"
)

/*
Env loads environment variables into the given EnvStruct either from the os environment or from the give EnvPath.

The EnvPath is optional. If no EnvPath is given, Env will attempt to load environment variables from the os environment.

The struct must have a tag named "barfenv" with the following format: "key=YOUR_ENV_KEY;required=true" or "key=YOUR_ENV_KEY;required=false".*/
func Env(EnvStruct interface{}, EnvPath ...string) error {
	if len(EnvPath) > 0 {
		// load environment variables
		if err := load(EnvPath[0]); err != nil {
			return err
		}
	}
	// ensure struct is not nil
	if EnvStruct == nil {
		return fmt.Errorf("EnvStruct must not be nil")
	}
	// verify argument is a pointer to a struct
	rv := reflect.ValueOf(EnvStruct)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct || rv.IsNil() {
		return fmt.Errorf("EnvStruct must be a pointer to a struct")
	}
	// verify environment variables
	if err := verify(rv); err != nil {
		return err
	}
	// prepare struct
	// load environment variables into struct
	append(rv)
	return nil
}
