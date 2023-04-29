/* package env
barf's simple interface for interacting environment variables. */
package env

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// load attempts to set key value pairs from the given file path
// to the os' environment.
// load also tries to avoid parsing comments and empty lines.
// load returns an error if the file path is invalid or if the
// file cannot be read.
func load(path string) error {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	// close file
	defer f.Close()
	// create scanner
	s := bufio.NewScanner(f)
	// scan file
	for s.Scan() {
		// get line
		line := s.Text()
		// skip empty lines
		if line == "" {
			continue
		}

		// handle comments

		if strings.Contains(line, "#") {

			// skip comments
			if strings.HasPrefix(line, "#") {
				continue
			}

			// use regex to extract key value pairs from line excluding comments
			if match, err := regexp.Compile(`^(?P<key>[^=]+)=(?P<value>[^#]+)`); err == nil {
				// check if line is valid
				if match.MatchString(line) {
					// get key value pairs
					pair := match.FindStringSubmatch(line)
					fmt.Println(pair)
					// set environment variable
					if err := os.Setenv(pair[1], pair[2]); err != nil {
						return err
					}
				}
				continue
			}
		}

		// split line
		pair := strings.Split(line, "=")
		// check if line is valid
		if len(pair) != 2 {
			return fmt.Errorf("invalid line in env file %s", line)
		}
		// set environment variable
		if err := os.Setenv(pair[0], pair[1]); err != nil {
			return err
		}
	}
	// check if scanner has an error
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}
