package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	db "gopkg.in/gorethink/gorethink.v3"
	"github.com/benbjohnson/clock"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// 
// main_test.go
// ContentService 
//
// Copyright © 2017 NGINX Inc. All rights reserved.
//

// Mock database for testing calls to RethinkDB
var mock = db.NewMock()
// Test environment to be injected into handlers when called
// Used to inject mock database into handlers for testing
var testEnv = &Env{
	Session: mock,
	Clock: clock.NewMock(),
}

func MockAlbumManager(t *testing.T)  {
	handler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request){
		fmt.Fprintf(resp, "hello, you've hit %s\n", req.URL.Path)
	})
	mockAlbumManager := httptest.NewServer(handler)
	mockAlbumManager.URL = os.Getenv("ALBUM_MANAGER_HOST")
	mockAlbumManager.Start()
	defer mockAlbumManager.Close()

}

// Tests NewArticle function
func TestNewArticle(t *testing.T) {
	// Specify values within Post object for replacing article
	id := `{"location":"locationCreate", "author":"nameCreate", "photo":"photoCreate", "title":"titleCreate", "extract":"extractCreate", "body":"bodyCreate", "album_id": 1}`
	post := Post{}
	json.Unmarshal([]byte(id), &post)
	post.Date = testEnv.Clock.Now()
	// Specify return variable for what should be returned by database
	var resp db.WriteResponse
	resp.Inserted = 1
	resp.GeneratedKeys = []string{"cc60e237-fa52-4b9c-9d72-de2ae808f535"}

	// Set database return values on reception of request to create article
	// with new Post object
	mock.On(db.DB("content").Table("posts").Insert(post)).Return(resp, nil)

	// Create new HTTP request
	// Method: POST
	// Pattern: /v1/content/
	// Body: New Post object
	req, err := http.NewRequest(http.MethodPost, "/v1/content", bytes.NewBufferString(id))
	if err != nil {
         t.Fatal(err)
	}
	req.Header.Add("auth-id", "12345678-abcd-1234-abcd-1234567890ab")

	// Initilize new recorder for testing response of handler
	rr := httptest.NewRecorder()

	// Call Articles handler
	NewRouter(testEnv).ServeHTTP(rr, req)

	// Check status code of response
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Compare expected value to body of response
	expected := `[cc60e237-fa52-4b9c-9d72-de2ae808f535]`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}
}

// Tests Articles function
func TestGetAllArticles(t *testing.T) {
	// Specify return variable for what should be returned by database
	var expected = []interface{}{
		map[string]interface{}{"author":"nameCreate","date":"2017-11-01T21:29:31.744Z","extract":"extractCreate","id":"cc60e237-fa52-4b9c-9d72-de2ae808f535","location":"locationCreate","photo":"photoCreate","title":"titleCreate","album_id": 1},
		map[string]interface{}{"author":"nameCreate","date":"2017-11-01T21:47:42.201Z","extract":"extractCreate","id":"4b8073ba-61d5-4626-a51c-992ceb6cd5d1","location":"locationCreate","photo":"photoCreate","title":"titleCreate","album_id": 1},
	}
	
	// Set database return values on reception of request to get all articles
	mock.On(db.DB("content").Table("posts").WithFields("id", "date", "location", "author", "photo", "title", "extract", "album_id")).Return(
		expected, nil)

	// Create new HTTP request
	// Method: GET
	// Pattern: /v1/content/
	// Body: None
	req, err := http.NewRequest(http.MethodGet, "/v1/content", nil)
	if err != nil {
         t.Fatal(err)
	}

	// Initilize new recorder for testing response of handler
	rr := httptest.NewRecorder()

	// Call GetAllArticles handler
	NewRouter(testEnv).ServeHTTP(rr, req)

	// Check status code of response
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}

	// Compare expected value to body of response (Not implemented)

}

// Test Article Function
func TestGetArticle(t *testing.T) {
	// Specify test article ID for which article to delete in database
	// Specify return variable for what should be returned by database
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`
	expected := map[string]interface{}{"author":"nameCreate","date":"2017-11-01T21:29:31.744Z","extract":"extractCreate","id":"cc60e237-fa52-4b9c-9d72-de2ae808f535", "location":"locationCreate","photo":"photoCreate","title":"titleCreate","album_id": 1}


	// Set database return values on reception of request to get article
	// with specified article ID
	mock.On(db.DB("content").Table("posts").Get(articleId).Pluck("id", "date",
		"location", "author", "photo", "title", "body", "extract", "album_id")).Return(
		expected, nil)

	// Create new HTTP request
	// Method: GET
	// Pattern: /v1/content/{articleId}
	// Body: None
	req, err := http.NewRequest(http.MethodGet, "/v1/content/" + articleId, nil)
	if err != nil {
         t.Fatal(err)
	}

	// Initilize new recorder for testing response of handler
	rr := httptest.NewRecorder()

	// Call GetArticle handler
	NewRouter(testEnv).ServeHTTP(rr, req)

	// Check status code of response
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}

	// Compare expected value to body of response (Not implemented)

}

// Test ReplaceArticle Function
func TestReplaceArticle(t *testing.T) {
	// Specify test article ID for which article to delete in database
	// Specify values within Post object for replacing article

	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`
	id := `{"author":"newAuthor"}`
	post := Post{}
	json.Unmarshal([]byte(id), &post)
	post.Date = testEnv.Clock.Now()
	post.Id = articleId

	var resp db.WriteResponse
	resp.Replaced = 2

	// Set database return values on reception of request to replace article
	// with specified article ID and Post object
	mock.On(db.DB("content").Table("posts").Get(articleId).Replace(post, db.ReplaceOpts{ReturnChanges: true})).Return(
		resp, nil)

	// Create new HTTP request
	// Method: PUT
	// Pattern: /v1/content/{articleId}
	// Body: new Post object
	req, err := http.NewRequest(http.MethodPut, "/v1/content/" + articleId, bytes.NewBufferString(id))
	if err != nil {
         t.Fatal(err)
	}
	req.Header.Add("auth-id", "12345678-abcd-1234-abcd-1234567890ab")

	// Initilize new recorder for testing response of handler
	rr := httptest.NewRecorder()

	// Call ReplaceArticle handler
	NewRouter(testEnv).ServeHTTP(rr, req)

	// Check status code of response
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Compare expected value to body of response
	expected := `{"Errors":0,"Inserted":0,"Updated":0,"Unchanged":0,"Replaced":2,"Renamed":0,"Skipped":0,"Deleted":0,"Created":0,"DBsCreated":0,"TablesCreated":0,"Dropped":0,"DBsDropped":0,"TablesDropped":0,"GeneratedKeys":null,"FirstError":"","ConfigChanges":null,"Changes":null}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}
}

// Test UpdateArticle function
func TestUpdateArticle(t *testing.T) {
	// Specify test article ID for which article to delete in database
	// Specify element + newValue pair for updating specified article
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`
	element := `author`
	newValue := `newValue`
	str := `{"` + element +`": "` + newValue + `"}`
	post := Post{}
	json.Unmarshal([]byte(str), &post)

	var resp db.WriteResponse
	resp.Replaced = 1

	// Set database return values on reception of request to update article
	// with specified article ID and element + newValue pair
	mock.On(db.DB("content").Table("posts").Get(articleId).Update(post, db.UpdateOpts{ReturnChanges: true})).Return(
		resp, nil)

	// Create new HTTP request
	// Method: PATCH
	// Pattern: /v1/content/{articleId}/{element}/{newValue}
	// Body: None
	req, err := http.NewRequest(http.MethodPatch, "/v1/content/" + articleId + "/" + element + "/" + newValue, nil)
	if err != nil {
         t.Fatal(err)
	}
	req.Header.Add("auth-id", "12345678-abcd-1234-abcd-1234567890ab")

	// Initilize new recorder for testing response of handler
	rr := httptest.NewRecorder()

	NewRouter(testEnv).ServeHTTP(rr, req)

	// Check status code of response
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Compare expected value to body of response
	expected := `{"Errors":0,"Inserted":0,"Updated":0,"Unchanged":0,"Replaced":1,"Renamed":0,"Skipped":0,"Deleted":0,"Created":0,"DBsCreated":0,"TablesCreated":0,"Dropped":0,"DBsDropped":0,"TablesDropped":0,"GeneratedKeys":null,"FirstError":"","ConfigChanges":null,"Changes":null}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}
}

// Test DeleteArticle function
func TestDeleteArticle(t *testing.T) {

	// Specify test article ID for which article to delete in database
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`

	// Set database return values on reception of request to delete article
	// with specified article ID
	mock.On(db.DB("content").Table("posts").Get(articleId).Delete(db.DeleteOpts{ReturnChanges: true})).Return(
		`{}`, nil)

	// Create new HTTP request
	// Method: DELETE
	// Pattern: /v1/content/{articleId}
	// Body: None
	req, err := http.NewRequest(http.MethodDelete, "/v1/content/" + articleId, nil)
	if err != nil {
         t.Fatal(err)
	}
	req.Header.Add("auth-id", "12345678-abcd-1234-abcd-1234567890ab")

	// Initilize new recorder for testing response of handler
	rr := httptest.NewRecorder()

	// Call DeleteArticle handler
	NewRouter(testEnv).ServeHTTP(rr, req)

	// Check status code of response
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}
	
	// Compare expected value to body of response
	expected :=`{}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}

	// Check if all database calls were made correctly
	mock.AssertExpectations(t)
}

// Test Album-Manager Set False/True
func TestAlbumManager(t *testing.T) {
	// Specify test article ID for which article to delete in database
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`
	albumID := 1
	id := `{"author":"newAuthor"}`

	// Create new HTTP request
	// Method: DELETE
	// Pattern: /v1/content/{articleId}
	// Body: None
	req, err := http.NewRequest(http.MethodPut, "/v1/content/" + articleId, bytes.NewBufferString(id))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("auth-id", "12345678-abcd-1234-abcd-1234567890ab")

	err = SetAlbumPublic(albumID, true, req)
	if err != nil {
		t.Errorf("SetAlbumPublic failed with error %s", err, "It should have worked")
	}
	err = SetAlbumPublic(albumID, false, req)
	if err != nil {
		t.Errorf("SetAlbumPublic failed with error %s", err, "It should have worked")
	}
}
