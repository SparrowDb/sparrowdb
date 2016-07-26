package spql

import "regexp"

var (
	// ValidatenDatabaseName Rule to accept only letters and numbers
	ValidateDatabaseName = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)
