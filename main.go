package main

import (
	db "gopkg.in/gorethink/gorethink.v3"
	"net/http"
	"log"
	"os"
	"fmt"
)

func main() {
	os.Setenv("RETHINKDB_URL", "localhost:28015")
	var session *db.Session
	var err error

	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := db.DBCreate("content").RunWrite(session)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("%d DB created", resp.DBsCreated)

	response, err := db.DB("content").TableCreate("posts").RunWrite(session)
	if err != nil {
		log.Print("Error creating table: %s", err)
	}

	fmt.Printf("%d table created", response.TablesCreated)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
