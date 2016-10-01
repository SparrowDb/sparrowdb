package errors

import "errors"

var (
	// ErrCreateDatabase error message when create database
	ErrCreateDatabase = errors.New("Could not create database")

	// ErrCreateSnapshot error message when create database snapshot
	ErrCreateSnapshot = errors.New("Could not create snapshot")

	// ErrDropDatabase error message when drop database
	ErrDropDatabase = errors.New("Could not drop database")

	// ErrOpenDatabase error message when open database
	ErrOpenDatabase = errors.New("Could not open database")

	// ErrDatabaseNotFound error message when don't find database
	ErrDatabaseNotFound = errors.New("Database not found")

	// ErrWrongRequest error message for wrong HTTP request
	ErrWrongRequest = errors.New("Wrong HTTP request")

	// ErrInvalidQueryAction error message for when query action is invalid
	ErrInvalidQueryAction = errors.New("Invalid query action")

	// ErrEmptyQueryResult error message for empty query result
	ErrEmptyQueryResult = errors.New("Empty query result")

	// ErrWrongToken error message when user inputs wrong token
	// for image request
	ErrWrongToken = errors.New("Wrong token")

	// ErrInsertImage error message when can't insert image
	ErrInsertImage = errors.New("Could not insert images")

	// ErrReadDir error message when try to read directory
	ErrReadDir = errors.New("Could not read directory")

	// ErrFileNotFound error message when file not found
	ErrFileNotFound = errors.New("File %s not found")

	// ErrParseFile error message when cannot parse file
	ErrParseFile = errors.New("Error trying to parse file %s")

	// ErrWrongInstanceMode error message for wrong instance mode
	ErrWrongInstanceMode = errors.New("Not valid SparrowDB mode, it must be [R]ead, [W]write or [RW]read-write")

	// ErrFileCorrupted error message when file is corrupted
	ErrFileCorrupted = errors.New("Could not read data from %s. File Corrupted")

	// ErrLogin error message when username and/or password is wrong
	ErrLogin = errors.New("Wrong username and/or password")

	// ErrInvalidToken error message when username inputs invalid or expired token
	ErrInvalidToken = errors.New("Invalid or expired token")
)
