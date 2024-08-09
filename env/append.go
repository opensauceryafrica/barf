/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/opensaucerer/barf/constant"
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
		tag := t.Type().Field(i).Tag.Get(constant.EnvTag)
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
			if _, ok := constant.EnvTagKeys[kv[0]]; !ok {
				continue
			}
			// set value
			barfenv[kv[0]] = kv[1]
		}
		// check if key is set
		if val, ok := barfenv["key"]; ok {
			// append environment variable to constant.Env
			switch t.Field(i).Kind() {
			case reflect.String:
				t.Field(i).SetString(Get(val))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				integer, err := strconv.ParseInt(Get(val), 10, 64)
				if err != nil {
					panic(err)
				}
				t.Field(i).SetInt(integer)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				integer, err := strconv.ParseUint(Get(val), 10, 64)
				if err != nil {
					panic(err)
				}
				t.Field(i).SetUint(integer)
			case reflect.Float32, reflect.Float64:
				float, err := strconv.ParseFloat(Get(val), 64)
				if err != nil {
					panic(err)
				}
				t.Field(i).SetFloat(float)
			case reflect.Bool:
				boolean, err := strconv.ParseBool(Get(val))
				if err != nil {
					panic(err)
				}
				t.Field(i).SetBool(boolean)
			}
		}
	}
}
