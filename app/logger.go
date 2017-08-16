package main

import (
	"log"
	"net/http"
	"time"
)

//
//  logger.go
//  ContentService
//
//  Copyright © 2017 NGINX Inc. All rights reserved.
//

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
