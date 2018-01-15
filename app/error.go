package main

//
// error.go
// ContentService
//
// Copyright © 2017 NGINX Inc. All rights reserved.

// Custom interface defines more easy to manage errors including a
// status function for returning the status code of a response
type Error interface {
	error
	Status() int
}

// Object for storing status code and error
type StatusError struct {
	Code int
	Err  error
}

// Function for defining what error to return
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Function for defining what status to return
func (se StatusError) Status() int {
	return se.Code
}
