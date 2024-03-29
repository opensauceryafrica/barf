package barf

import "github.com/opensaucerer/barf/env"

/*
Env loads environment variables into the given EnvStruct either from the os environment or from the give EnvPath.

The EnvPath is optional. If no EnvPath is given, Env will attempt to load environment variables from the os environment.

The struct provided must have a tag named "barfenv" with the following format:

	`barfenv:"key=YOUR_ENV_KEY"`
	`barfenv:"key=YOUR_ENV_KEY;required=true"`
	`barfenv:"key=YOUR_ENV_KEY;required=false"`

	`required` defaults to false if not specified.
*/
var Env = env.Env
