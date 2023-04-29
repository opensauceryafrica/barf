/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"os"
)

// Get retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
func Get(key string) string {
	return os.Getenv(key)
}
