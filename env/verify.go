/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/opensaucerer/barf/constant"
	"github.com/opensaucerer/barf/typing"
)

// verify checks if all environment variables are set correctly
// it uses the struct tag "barfenv" to verify the environment variables
func verify(env reflect.Value) error {
	// get the type of argument
	t := env.Elem()
	// verify each struct field tag
	for i := 0; i < t.NumField(); i++ {
		// get the field tag value
		tag := t.Type().Field(i).Tag.Get(constant.EnvTag)
		if tag == "" {
			continue
		}
		// default value is required
		required := false
		barfenv := typing.M{}
		// extract all key value pairs
		for _, pair := range strings.Split(tag, ";") {
			// split key value pair
			kv := strings.Split(pair, "=")
			// check if key value pair is valid
			if len(kv) != 2 {
				return fmt.Errorf("invalid key value pair -> %s <- in env struct %T", pair, env)
			}
			// check if key is valid
			val, ok := constant.EnvTagKeys[kv[0]]
			if !ok {
				return fmt.Errorf("invalid key %s", kv[0])
			}
			// check if value is valid
			if val != nil {
				// check if value is valid
				found := false
				switch reflect.TypeOf(val).Kind() {
				case reflect.Slice:
					// check if value is in slice
					for _, v := range reflect.ValueOf(val).Interface().([]string) {
						if v == kv[1] {
							found = true
							break
						}
					}

				case reflect.Map:
					// check if value is in map
					if _, ok := reflect.ValueOf(val).Interface().(map[string]interface{})[kv[1]]; ok {
						found = true
					}

				case reflect.String:
					// check if value is equal to string
					if reflect.ValueOf(val).Interface().(string) == kv[1] {
						found = true
					}

				default:
					return fmt.Errorf("invalid value type %T", val)

				}
				// return error if value is not valid
				if !found {
					return fmt.Errorf("invalid value %s for key %s", kv[1], kv[0])
				}
			}
			// set value
			barfenv[kv[0]] = kv[1]
		}
		// check if required is set
		if val, ok := barfenv["required"]; ok {
			// set required
			required = val == "true"
		}
		// check if key is set
		if val, ok := barfenv["key"]; ok {
			// check if required
			if required {
				// check if environment variable is set
				if Get(val) == "" {
					return fmt.Errorf("environment variable %s is not set", val)
				}
			}
		}

	}
	return nil
}
