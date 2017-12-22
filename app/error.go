package main

//
//  error.go
//  ContentService
//
//  Copyright © 2017 NGINX Inc. All rights reserved.
//

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}
