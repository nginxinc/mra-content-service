package main

import (
	"github.com/gorilla/mux"
)

//
// router.go
// ContentService
//
// Copyright © 2017 NGINX Inc. All rights reserved.
//

// Router for associating HTTP requests with functions based on URI
// Router takes parameters: Name, Method, Path, and Handler to associate
// with a function
func NewRouter(env *Env) *mux.Router {

	// Create new gorilla/mux router with with strict slash
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		// Associate each route with an HTTP endpoint
		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(Handler{env, route.Function})

	}

	// Return router to be used by server
	return router
}
