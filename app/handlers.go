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
//  handlers.go
//  ContentService
//
//  Copyright © 2017 NGINX Inc. All rights reserved.
//

type Photo struct {
	Name	string	`gorethink:"name"`
	Url		string	`gorethink:"url"`
}

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

type Env struct {
	Session  db.QueryExecutor
	IsTest	 bool
}

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

type HandlerFunc func(e *Env, w http.ResponseWriter, r *http.Request) error

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

func Welcome(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Welcome to the content service!")
	return nil
}

func Articles(env *Env, w http.ResponseWriter, r *http.Request) error {
	var resp *db.Cursor
	var err error

	resp, err = db.DB("content").Table("posts").WithFields("id", "date", "location", "author", "photo", "title", "extract").Run(env.Session)
	if err != nil {
		fmt.Print(err)
		return StatusError{500, err}
	}
	defer resp.Close()

	var posts []interface{}
	err = resp.All(&posts)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return StatusError{500, err}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		panic(err)
	}

	fmt.Printf("%d posts", len(posts))
	return nil
}

func Article(env *Env, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var articleId string
	var resp *db.Cursor
	var err error

	articleId = vars["articleId"]

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
	err = resp.All(&post)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return StatusError{500, err}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		panic(err)
	}
	return nil
}

func NewArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var newPost Post
	var resp db.WriteResponse
	err := decoder.Decode(&newPost)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

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

	fmt.Printf("%d row inserted, %d key generated", resp.Inserted, len(resp.GeneratedKeys))
	fmt.Fprint(w, resp.GeneratedKeys)

	return nil
}

func ReplaceArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var articleId string
	articleId = vars["articleId"]
	var resp db.WriteResponse

	decoder := json.NewDecoder(r.Body)
	var newPost Post
	err := decoder.Decode(&newPost)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
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

	printObj(w, resp)

	return nil
}

func UpdateArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var articleId string = vars["articleId"]
	var element string = vars["element"]
	var newValue string = vars["newValue"]

	var resp db.WriteResponse
	var err error

	str := `{"` + element +`": "` + newValue + `"}`
	fmt.Print(str)
	res := Post{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)

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

	printObj(w, resp)

	return nil
}

func DeleteArticle(env *Env, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var articleId string = vars["articleId"]

	var resp *db.Cursor
	var err error

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

	printObj(w, resp)

	return nil
}

func printObj(w http.ResponseWriter, v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Fprint(w, string(vBytes))
}
