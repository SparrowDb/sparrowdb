package spql

import "regexp"

const (
	errDataNotFound = "Data %s not found in %s"
)

var (
	// ValidateDatabaseName Rule to accept only letters and numbers
	ValidateDatabaseName = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)

// ValidateFieldName validates fild name
func ValidateFieldName(name string) bool {
	return name == "key" || name == "extension"
}

// ValidateMatchType validates match type of select query
func ValidateMatchType(name string) bool {
	return name == "*" || name == "="
}
