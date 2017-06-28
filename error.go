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
