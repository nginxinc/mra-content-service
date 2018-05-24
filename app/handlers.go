package main

import (
	"fmt"
	"net/http"
	"time"
	db "gopkg.in/gorethink/gorethink.v3"
	"github.com/benbjohnson/clock"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
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
	Album_id  int   	`gorethink:"album_id,omitempty" json:"album_id"`
}

// Environment object used to inject database state into handlers
// Practically used to test handlers with a mock database
type Env struct {
	Session  db.QueryExecutor
	Clock	 clock.Clock
}

// Handler object used for allowing handler functions to accept
// an environment object
type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}
// HandlerFunc type used to specify a template for handler functions to follow
type HandlerFunc func(e *Env, w http.ResponseWriter, r *http.Request) error

// ServeHTTP is called on each HTTP request. Specifies which function is
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
	resp, err = db.DB("content").Table("posts").WithFields(
		"id", "date", "location", "author", "photo", "title", "extract", "album_id").Run(env.Session)
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

	// Make call to rethink database
	resp, err = db.DB("content").Table("posts").Get(articleId).Pluck("id",
		"date", "location", "author", "photo", "title", "body", "extract", "album_id").Run(env.Session)
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
		log.Printf("HTTP %d - %s", err, err)
		return StatusError{500, err}
	}
	defer r.Body.Close()

	newPost.Date = env.Clock.Now()

	err = SetAlbumPublic(newPost.Album_id, true, r)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	// Make call to rethink database
	resp, err = db.DB("content").Table("posts").Insert(newPost).RunWrite(env.Session)
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
// @return: JSON object that specified what information was changed
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

	// Set id of article to fetch
	newPost.Id = articleId
	newPost.Date = env.Clock.Now()

	// Make call to rethink database and get changes back
	resp, err = db.DB("content").Table("posts").Get(articleId).Replace(newPost, db.ReplaceOpts{ReturnChanges: true}).RunWrite(env.Session)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	// Coming through as interface{}
	newAlbumID := getAlbumIDFromReturnValues(resp.Changes[0].NewValue)
	oldAlbumID := getAlbumIDFromReturnValues(resp.Changes[0].OldValue)

	if newAlbumID != oldAlbumID {
		err = SetAlbumPublic(newAlbumID, true, r)
		if err != nil {
			fmt.Print(err)
			return StatusError{500, err}
		}
		err = SetAlbumPublic(oldAlbumID, false, r)
		if err != nil {
			fmt.Print(err)
			return StatusError{500, err}
		}
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
// @return: JSON object that specified what information was changed
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

	// Make call to rethink database and get changes back
	resp, err = db.DB("content").Table("posts").Get(articleId).Update(res, db.UpdateOpts{ReturnChanges: true}).RunWrite(env.Session)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	// Coming through as interface{}
	newAlbumID := getAlbumIDFromReturnValues(resp.Changes[0].NewValue)
	oldAlbumID := getAlbumIDFromReturnValues(resp.Changes[0].OldValue)

	if newAlbumID != oldAlbumID {
		err = SetAlbumPublic(newAlbumID, true, r)
		if err != nil {
			fmt.Print(err)
			return StatusError{500, err}
		}
		err = SetAlbumPublic(oldAlbumID, false, r)
		if err != nil {
			fmt.Print(err)
			return StatusError{500, err}
		}
	}

	// Print JSON response
	printObj(w, resp)

	return err
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

	// Make call to rethink database
	resp, err = db.DB("content").Table("posts").Get(articleId).Delete(db.DeleteOpts{ReturnChanges: true}).Run(env.Session)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	result , _ := resp.NextResponse()

	// To get the albumID out of the response takes 5 levels of indirection
	// result->deRefedJson->jsonResult->changes->change->albumID
	// unfortunately, a Struct to extract the value didn't seem to work because of the
	// anonymous []interface{} in the jsonResult
	var deRefedJson *json.RawMessage
	var jsonResult map[string]interface{}
	var changes map[string]interface{}
	var change map[string]interface{}

	deRefedJson = (*json.RawMessage)(&result)
	err = json.Unmarshal(*deRefedJson, &jsonResult)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}

	changesInterfaces := jsonResult["changes"].([]interface{})
	changes = changesInterfaces[0].(map[string]interface{})
	change = changes["old_val"].(map[string]interface{})
	albumID := change["album_id"].(float64)

	err = SetAlbumPublic(int(albumID), false, r)
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

func getAlbumIDFromReturnValues(data interface{}) int {
	m := data.(map[string]interface{})
	var albumID int
	if albumFloat, ok := m["album_id"].(float64); ok {
		albumID = int(albumFloat)
	}
	return albumID
}
