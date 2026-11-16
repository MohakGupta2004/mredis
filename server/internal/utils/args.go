package utils

// Args look like: [key, value1, value2, ...]

// GetKey returns the key, which is always the first arg.
func GetKey(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

// GetValues returns every value after the key.
func GetValues(args []string) []string {
	if len(args) < 2 {
		return []string{}
	}
	return args[1:]
}

// GetValue returns one value after the key. index 0 is the first value after the key.
// Returns "" when that value was not sent.
func GetValue(args []string, index int) string {
	values := GetValues(args)
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}
