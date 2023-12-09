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
	// if value is nil, return default value
	if value == nil {
		return def
	}

	// if value is not nil, check inner value type
	switch v := value.(type) {
	case string:
		if v == "" {
			return def
		}
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float64, float32, bool:
		return v

	}

	// return value if not nil
	return value

}
