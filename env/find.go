/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"os"
)

// find searches for the given key in the environment variables
// and returns the value and a boolean indicating if the key was found
func Find(key string) (string, bool) {
	return os.LookupEnv(key)
}
