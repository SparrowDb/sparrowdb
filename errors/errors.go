package errors

import "errors"

var (
	// ErrCreateDatabase error message when create database
	ErrCreateDatabase = errors.New("Could not create database")

	// ErrDropDatabase error message when drop database
	ErrDropDatabase = errors.New("Could not drop database")

	// ErrOpenDatabase error message when open database
	ErrOpenDatabase = errors.New("Could not open database")

	// ErrDatabaseNotFound error message when don't find database
	ErrDatabaseNotFound = errors.New("Database not found")

	// ErrWrongRequest error message for wrong HTTP request
	ErrWrongRequest = errors.New("Wrong HTTP request")

	// ErrEmptyQueryResult error message for empty query result
	ErrEmptyQueryResult = errors.New("Empty query result")

	// ErrWrongToken error message when user inputs wrong token
	// for image request
	ErrWrongToken = errors.New("Wrong token")

	// ErrInsertImage error message when can't insert image
	ErrInsertImage = errors.New("Could not insert images")

	// ErrReadDir error message when try to read directory
	ErrReadDir = errors.New("Could not read directory")

	// ErrFileCorrupted erros message when file is corrupted
	ErrFileCorrupted = errors.New("Could not read data from %s. File Corrupted")
)
