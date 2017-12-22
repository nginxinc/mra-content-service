package main

import (
	// "net/http"

	"github.com/gorilla/mux"
)

//
//  router.go
//  ContentService
//
//  Copyright © 2017 NGINX Inc. All rights reserved.
//

func NewRouter(env *Env) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		// var handler http.Handler

		// handler = route.HandlerFunc
		// handler = Logger(handler, route.Name)

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(Handler{env, route.Function})

	}

	return router
}
