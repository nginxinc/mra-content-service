package main

import (
	//"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"strconv"

	//"github.com/gorilla/mux"

)

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the content service!")
}

func Articles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Articles")
}

func Article(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Article")
}

func NewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "New article")
}

func ReplaceArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Replace article")
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Update article")
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Delete article")
}
