package main

import (
	"log"
	"net/http"
	"time"
	"os"
)

//
// logger.go
// ContentService
//
// Copyright © 2017 NGINX Inc. All rights reserved.
//

// Function for logging incoming request information
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.SetOutput(os.Stdout)

		inner.ServeHTTP(w, r)

		// Log request method, URL, name of function, and current time
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
