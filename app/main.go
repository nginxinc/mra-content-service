package main

import (
	db "gopkg.in/gorethink/gorethink.v3"
	"github.com/joho/godotenv"
	"net/http"
	"log"
	"os"
	"fmt"
)

//
//  main.go
//  ContentService
//
//  Copyright © 2017 NGINX Inc. All rights reserved.
//

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var session *db.Session

	fmt.Print(os.Getenv("RETHINKDB_URL"))
	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	env := &Env{
		Session: session,
		IsTest: false,
    }

	resp, err := db.DBCreate("content").RunWrite(env.Session)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("%d DB created", resp.DBsCreated)

	response, err := db.DB("content").TableCreate("posts").RunWrite(env.Session)
	if err != nil {
		log.Print("Error creating table: " + err.Error())
	}

	fmt.Printf("%d table created", response.TablesCreated)

	router := NewRouter(env)

	log.Fatal(http.ListenAndServe(":8080", router))
}
