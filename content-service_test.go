package main

import (
	"testing"
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
	"time"
	"encoding/json"
)

var urlService string = "http://content-service.mra.nginxps.com/"
var client = &http.Client{
	Timeout: 5 * time.Second,
}

func createArticle(bodyStr string) string {
	req, err := http.NewRequest("POST", urlService + "v1/content", bytes.NewBufferString(bodyStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body[1:len(body)-1])
}

func replaceArticle(id string, bodyStr string) string {
	req, err := http.NewRequest("PUT", urlService + "v1/content/" + id, bytes.NewBufferString(bodyStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body[1:len(body)-2])
}

func deleteArticle(id string) {
	req, err := http.NewRequest("DELETE", urlService + "v1/content/" + id, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func getArticle(id string) string {
	req, err := http.NewRequest("GET", urlService + "v1/content/" + id, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if string(body) != "" {
		return string(body[1:len(body)-2])
	} else {
		return string(body)
	}
}

func TestWelcome(t *testing.T) {
	resp, _ := http.Get(urlService)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	newStr := buf.String()

	if newStr != "Welcome to the content service!" {
		t.Error()
	}
}

func TestNewArticle(t *testing.T) {
	id := createArticle(`{"location":"locationCreate", "author":"nameCreate", "photo":"photoCreate", "title":"titleCreate", "extract":"extractCreate", "body":"bodyCreate"}`)
	body := getArticle(id)

	// Compare bodies
	article := Post{}
	json.Unmarshal([]byte(body), &article)
	if (article.Author != "nameCreate" || article.Body != "bodyCreate" || article.Id != id ||
		article.Location != "locationCreate" || article.Photo != "photoCreate" || article.Title != "titleCreate") {
		t.Error()
	}
}

func TestReplaceArticle(t *testing.T) {
	id := createArticle(`{"location":"locationCreate", "author":"nameCreate", "photo":"photoCreate", "title":"titleCreate", "extract":"extractCreate", "body":"bodyCreate"}`)
	bodyStr := `{"location":"locationReplace", "author":"nameReplace", "photo":"photoReplace", "title":"titleReplace", "extract":"extractReplace", "body":"bodyReplace"}`
	replaceArticle(id, bodyStr)
	body := getArticle(id)

	// Compare bodies
	article := Post{}
	json.Unmarshal([]byte(body), &article)
	if (article.Author != "nameReplace" || article.Body != "bodyReplace" || article.Id != id ||
		article.Location != "locationReplace" || article.Photo != "photoReplace" || article.Title != "titleReplace") {
		t.Error()
	}
}

func TestDeleteArticle(t *testing.T) {
	id := createArticle(`{"location":"locationCreate", "author":"nameCreate", "photo":"photoCreate", "title":"titleCreate", "extract":"extractCreate", "body":"bodyCreate"}`)
	deleteArticle(id)
	body := getArticle(id)
	if body != "" {
		t.Error()
	}
}
