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
