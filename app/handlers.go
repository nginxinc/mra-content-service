package main

import (
	"fmt"
	"net/http"
	"time"
	db "gopkg.in/gorethink/gorethink.v3"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"strings"
)

//
// handlers.go
// ContentService
//
// Copyright © 2017 NGINX Inc. All rights reserved.
//

// Struct object for storing photo information within RethinkDB
type Photo struct {
	Name	string	`gorethink:"name"`
	Url		string	`gorethink:"url"`
}

// Struct object for storing post information within RethinkDB
type Post struct {
	Id		  string	`gorethink:"id,omitempty"`
	Date      time.Time	`gorethink:"date" json:"date"`
	Location  string	`gorethink:"location,omitempty" json:"location"`
	Author    string	`gorethink:"author,omitempty" json:"author"`
	Photo     string	`gorethink:"photo,omitempty" json:"photo"`
	Title     string	`gorethink:"title,omitempty" json:"title"`
	Extract   string	`gorethink:"extract,omitempty" json:"extract"`
	Body      string	`gorethink:"body,omitempty" json:"body"`
}

// Environment object used to inject database state into handlers
// Practically used to test handlers with a mock database
// (Should remove IsTest in future releases when mux is fixed)
type Env struct {
	Session  db.QueryExecutor
	IsTest	 bool
}

// Handler object used for allowing handler functions to accept
// an environment object
type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}
// HandlerFunc type used to specify a template for handler functions to follow
type HandlerFunc func(e *Env, w http.ResponseWriter, r *http.Request) error

// ServeHTTP is called on each HTTP request. Species which function is
// called as well as how errors are handled and how logging is set
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

// Handler listening for GET at "/" URI
// Returns specified string
// @return: string "Welcome to the content service!"
func Welcome(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Welcome to the content service!")
	return nil
}

// Handler listening for GET at "/v1/content" URI
// Get array of all articles in database
// @return: array of JSON objects, each element with single post information
func GetAllArticles(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp *db.Cursor
	var err error

	// Call database and get all articles with fields:
	// id, date, location, author, photo, title, and extract
	resp, err = db.DB("content").Table("posts").WithFields("id", "date", "location", "author", "photo", "title", "extract").Run(env.Session)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}
	defer resp.Close()

	var posts []interface{}
	// Map all posts into an array of interfaces for printing to user
	err = resp.All(&posts)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return StatusError{500, err}
	}

	// Write header to specify returning a JSON object
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		panic(err)
	}

	fmt.Printf("%d posts", len(posts))
	return nil
}

// Handler listening for GET at "/v1/content/{articleID}" URI
// Get specified article
// Parameters: articleID - specifies which article to get
// @return: single JSON object with post information
func GetArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp *db.Cursor
	var err error

	// Read gorilla/mux variables for article ID to fetch from database
	vars := mux.Vars(r)
	var articleId string = vars["articleId"]

	// IsTest variable shouldn't be necessary, but setting mux variables
	// must be set correctly in tests in order to avoid this
	if env.IsTest {
		resp, err = db.DB("content").Table("posts").Get(strings.Split(r.URL.Path, "/")[3]).Pluck("id", "date", "location", "author", "photo", "title", "body").Run(env.Session)
	} else {
		resp, err = db.DB("content").Table("posts").Get(articleId).Pluck("id", "date", "location", "author", "photo", "title", "body").Run(env.Session)
	}
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	defer resp.Close()

	var post []interface{}
	// Map all posts into an array of interfaces for printing to user
	err = resp.All(&post)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return StatusError{500, err}
	}

	// Write header to specify returning a JSON object
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		panic(err)
	}
	return nil
}

// Handler listening for POST at "/v1/content" URI
// Creates new article based on JSON object in POST
// @POST: new post object
// @return: ID of post within database
func NewArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp db.WriteResponse

	// Decode request body into Post object
	var newPost Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPost)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// IsTest variable shouldn't be necessary, but setting mux variables
	// must be set correctly in tests in order to avoid this
	if env.IsTest {
		resp, err = db.DB("content").Table("posts").Insert(newPost).RunWrite(env.Session)
	} else {
		newPost.Date = time.Now()
		resp, err = db.DB("content").Table("posts").Insert(newPost).RunWrite(env.Session)
	}
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	// Return number of rows inserted into database and keys generated
	fmt.Printf("%d row inserted, %d key generated", resp.Inserted, len(resp.GeneratedKeys))
	fmt.Fprint(w, resp.GeneratedKeys)

	return nil
}

// Handler listening for PUT at "/v1/content/{articleID}" URI
// Updates elements within specified post
// Parameters: articleID - specifies which article to update
// @POST: new post object
// @return: JSON object that specied what information was changed
func ReplaceArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp db.WriteResponse

	// Read gorilla/mux variables for article ID to fetch from database
	vars := mux.Vars(r)
	var articleId string = vars["articleId"]

	// Decode request body into Post object
	var newPost Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPost)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	// IsTest variable shouldn't be necessary, but setting mux variables
	// must be set correctly in tests in order to avoid this
	if env.IsTest {
		resp, err = db.DB("content").Table("posts").Get(strings.Split(r.URL.Path, "/")[3]).Replace(newPost).RunWrite(env.Session)
	} else {
		newPost.Id = articleId
		newPost.Date = time.Now()
		resp, err = db.DB("content").Table("posts").Get(articleId).Replace(newPost).RunWrite(env.Session)
	}
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	// Print JSON response
	printObj(w, resp)

	return nil
}

// Handler listening for PATCH at "/v1/content/{articleId}/{element}/{newValue}" URI
// Updates single element within specified post
// Parameters: articleID - specifies which article to update
// 			   element - element within post to update
// 			   newValue - new value of element
// @return: JSON object that specied what information was changed
func UpdateArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp db.WriteResponse
	var err error

	// Read gorilla/mux variables for article ID and element + value to update database
	vars := mux.Vars(r)
	var articleId string = vars["articleId"]
	var element string = vars["element"]
	var newValue string = vars["newValue"]

	// Set syntax for new/updated element within database
	str := `{"` + element +`": "` + newValue + `"}`

	// Unmarshal str variable into Post object
	res := Post{}
	json.Unmarshal([]byte(str), &res)

	// IsTest variable shouldn't be necessary, but setting mux variables
	// must be set correctly in tests in order to avoid this
	if env.IsTest {
		s := strings.Split(r.URL.Path, "/")
		resp, err = db.DB("content").Table("posts").Get(s[3]).Update(`{"` + s[4] +`": "` + s[5] + `"}`).RunWrite(env.Session)
	} else {
		resp, err = db.DB("content").Table("posts").Get(articleId).Update(res).RunWrite(env.Session)
	}
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	// Print JSON response
	printObj(w, resp)

	return nil
}

// Handler listening for DELETE at "/v1/content/{articleId}" URI
// Delete specified post
// Parameters: articleID - specifies which article to delete
// @return: empty JSON object
func DeleteArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp *db.Cursor
	var err error

	// Read gorilla/mux variables for article ID and element + value to update database
	vars := mux.Vars(r)
	var articleId string = vars["articleId"]

	// IsTest variable shouldn't be necessary, but setting mux variables
	// must be set correctly in tests in order to avoid this
	if env.IsTest {
		resp, err = db.DB("content").Table("posts").Get(strings.Split(r.URL.Path, "/")[3]).Delete().Run(env.Session)
	} else {
		resp, err = db.DB("content").Table("posts").Get(articleId).Delete().Run(env.Session)
	}
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	defer resp.Close()

	// Print JSON response
	printObj(w, resp)

	return nil
}

// Print JSON object to screen
func printObj(w http.ResponseWriter, v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Fprint(w, string(vBytes))
}
