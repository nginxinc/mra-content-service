package main

import (
	"fmt"
	"net/http"
	"time"
	db "gopkg.in/gorethink/gorethink.v3"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
)

type Photo struct {
	Name    string `gorethink:"name"`
	Url   string `gorethink:"url"`
}

type Post struct {
	Id	string		`gorethink:"id,omitempty"`
	Date     time.Time         `gorethink:"date" json:"date"`
	Location  string       `gorethink:"location,omitempty" json:"location"`
	Author  string       `gorethink:"author,omitempty" json:"author"`
	Photo  string       `gorethink:"photo,omitempty" json:"photo"`
	Title  string       `gorethink:"title,omitempty" json:"title"`
	Extract  string       `gorethink:"extract,omitempty" json:"extract"`
	Body  string       `gorethink:"body,omitempty" json:"body"`
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the content service!")
}

func Articles(w http.ResponseWriter, r *http.Request) {
	fmt.Print(os.Getenv("RETHINKDB_URL"))
	var session *db.Session
	var err error

	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := db.DB("content").Table("posts").WithFields("id", "date", "location", "author", "photo", "title", "extract").Run(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Close()

	var posts []interface{}
	err = resp.All(&posts)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		panic(err)
	}

	fmt.Printf("%d posts", len(posts))
}

func Article(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var articleId string

	articleId = vars["articleId"]

	var session *db.Session
	var err error

	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := db.DB("content").Table("posts").Get(articleId).Pluck("id", "date", "location", "author", "photo", "title", "body").Run(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Close()

	var post []interface{}
	err = resp.All(&post)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		panic(err)
	}
}

func NewArticle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newPost Post
	err := decoder.Decode(&newPost)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var session *db.Session
	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	newPost.Date = time.Now()
	resp, err := db.DB("content").Table("posts").Insert(newPost).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row inserted, %d key generated", resp.Inserted, len(resp.GeneratedKeys))
	fmt.Fprint(w, resp.GeneratedKeys)
}

func ReplaceArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var articleId string
	articleId = vars["articleId"]

	decoder := json.NewDecoder(r.Body)
	var newPost Post
	err := decoder.Decode(&newPost)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	newPost.Id = articleId

	var session *db.Session
	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	newPost.Date = time.Now()
	resp, err := db.DB("content").Table("posts").Get(articleId).Replace(newPost).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	printObj(resp)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var articleId string
	articleId = vars["articleId"]
	var element string
	var newValue string
	element = vars["element"]
	newValue = vars["newValue"]

	str := `{"` + element +`": "` + newValue + `"}`
	fmt.Print(str)
	res := Post{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)

	var session *db.Session
	var err error
	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := db.DB("content").Table("posts").Get(articleId).Update(res).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	printObj(resp)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var articleId string
	articleId = vars["articleId"]

	var session *db.Session
	var err error

	session, err = db.Connect(db.ConnectOpts{
		Address: os.Getenv("RETHINKDB_URL"),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := db.DB("content").Table("posts").Get(articleId).Delete().Run(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Close()

	fmt.Print("1 deleted article")
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}
