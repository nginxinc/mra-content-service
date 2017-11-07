package main

import (
	db "gopkg.in/gorethink/gorethink.v3"
	"net/http"
	"log"
	"os"
	"fmt"
)

//
// main.go
// ContentService
//
// Copyright © 2017 NGINX Inc. All rights reserved.
//

// Main function responsible for calling initializer functions for
// establishing HTTP endpoints with app listening on port :8080
func main() {

	fmt.Print(os.Getenv("RETHINKDB_URL"))

	// Initialize session variable by connecting to RethinkDB database
	// Specified by "RETHINKDB_URL" environment variable
	session, err := db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Initialize environment variable to inject into handlers.
	// IsTest set to false because this is a production environment (will remove in later release)
	env := &Env{
		Session: session,
		IsTest: false,
    }

	// Create database called "content" for storing articles
	resp, err := db.DBCreate("content").RunWrite(env.Session)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("%d DB created", resp.DBsCreated)

	// Create table called "posts" within database
	response, err := db.DB("content").TableCreate("posts").RunWrite(env.Session)
	if err != nil {
		log.Print("Error creating table: " + err.Error())
	}

	fmt.Printf("%d table created", response.TablesCreated)

	// Initialize router for mapping functions within handlers.go
	// to HTTP endpoints at specified URIs within routes.go
	router := NewRouter(env)

	// Listen for requests on port :8080 with router and logging
	log.Fatal(http.ListenAndServe(":8080", router))
}
