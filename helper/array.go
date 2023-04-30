package helper

type StringArray []string

// Contains returns true if the given string is in the array.
func (a StringArray) Contains(s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}

// Index returns the index of the given string in the array.
func (a StringArray) Index(s string) int {
	for i, v := range a {
		if v == s {
			return i
		}
	}
	return -1
}

// Remove removes the first occurrence of the given string from the array.
func (a StringArray) Remove(s string) StringArray {
	index := a.Index(s)
	if index >= 0 {
		return append(a[:index], a[index+1:]...)
	}
	return a
}

// Add adds the given string to the array.
func (a StringArray) Add(s string) StringArray {
	return append(a, s)
}

// AddUnique adds the given string to the array if it doesn't exist.
func (a StringArray) AddUnique(s string) StringArray {
	if !a.Contains(s) {
		return append(a, s)
	}
	return a
}

// AddAll adds all the given strings to the array.
func (a StringArray) AddAll(s []string) StringArray {
	for _, v := range s {
		a = a.Add(v)
	}
	return a
}

// Transform transforms the array using the given function.
func (a StringArray) Transform(f func(string) string) StringArray {
	for i, v := range a {
		a[i] = f(v)
	}
	return a
}
