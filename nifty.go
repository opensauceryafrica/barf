package barf

var (
	// allow is the true value
	allow = true

	// disallow is the false value
	disallow = false
)

// Allow is a nifty function that returns a pointer to a bool with value true.
func Allow() *bool {
	return &allow
}

// Disallow is a nifty function that returns a pointer to a bool with value false.
func Disallow() *bool {
	return &disallow
}

// True is a nifty function that returns a bool with value true.
func True() bool {
	return allow
}

// False is a nifty function that returns a bool with value false.
func False() bool {
	return disallow
}

// Obtain is a nifty function that returns the value passed to it or the default value passed, if the value passed is nil.
func Obtain(value, def interface{}) interface{} {
	if value == nil {
		return def
	}
	return value
}
