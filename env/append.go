/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"reflect"
	"strings"

	"github.com/opensaucerer/barf/config"
	"github.com/opensaucerer/barf/typing"
)

// append loads the environment variables into the provided struct
// the struct must have a tag named "barfenv" with the following format: "key=value;key=value;..."
func append(env reflect.Value) {
	// get the type of argument
	t := env.Elem()
	// append each struct field tag
	for i := 0; i < t.NumField(); i++ {
		// get the field tag value
		tag := t.Type().Field(i).Tag.Get(config.EnvTag)
		if tag == "" {
			continue
		}
		barfenv := typing.M{}
		// extract all key value pairs
		for _, pair := range strings.Split(tag, ";") {
			// split key value pair
			kv := strings.Split(pair, "=")
			// check if key value pair is valid
			if len(kv) != 2 {
				continue
			}
			// check if key is valid
			if _, ok := config.EnvTagKeys[kv[0]]; !ok {
				continue
			}
			// set value
			barfenv[kv[0]] = kv[1]
		}
		// check if key is set
		if val, ok := barfenv["key"]; ok {
			// append environment variable to constant.Env
			t.Field(i).SetString(Get(val))
		}
	}
}
