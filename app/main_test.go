package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	db "gopkg.in/gorethink/gorethink.v3"
	"bytes"
	"encoding/json"
)

var mock = db.NewMock()
var testEnv = &Env{
	Session: mock,
	IsTest: true,
}


func TestNewArticle(t *testing.T) {
	id := `{"location":"locationCreate", "author":"nameCreate", "photo":"photoCreate", "title":"titleCreate", "extract":"extractCreate", "body":"bodyCreate"}`
	post := Post{}
	json.Unmarshal([]byte(id), &post)

	var resp db.WriteResponse
	resp.Inserted = 1
	resp.GeneratedKeys = []string{"cc60e237-fa52-4b9c-9d72-de2ae808f535"}

	expected := `[cc60e237-fa52-4b9c-9d72-de2ae808f535]`
	mock.On(db.DB("content").Table("posts").Insert(post)).Return(
		resp, nil)

	req, err := http.NewRequest("POST", "/v1/content", bytes.NewBufferString(id))
	if err != nil {
         t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := Handler{testEnv, NewArticle}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Compare bodies
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}
}

func TestArticles(t *testing.T) {
	var expected = []interface{}{
		map[string]interface{}{"author":"nameCreate","date":"2017-11-01T21:29:31.744Z","extract":"extractCreate","id":"cc60e237-fa52-4b9c-9d72-de2ae808f535","location":"locationCreate","photo":"photoCreate","title":"titleCreate"},
		map[string]interface{}{"author":"nameCreate","date":"2017-11-01T21:47:42.201Z","extract":"extractCreate","id":"4b8073ba-61d5-4626-a51c-992ceb6cd5d1","location":"locationCreate","photo":"photoCreate","title":"titleCreate"},
	}
	
	mock.On(db.DB("content").Table("posts").WithFields("id", "date", "location", "author", "photo", "title", "extract")).Return(
		expected, nil)

	req, err := http.NewRequest("GET", "/v1/content", nil)
	if err != nil {
         t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := Handler{testEnv, Articles}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}

func TestArticle(t *testing.T) {
	var expected = map[string]interface{}{"author":"nameCreate","date":"2017-11-01T21:29:31.744Z","extract":"extractCreate","id":"cc60e237-fa52-4b9c-9d72-de2ae808f535","location":"locationCreate","photo":"photoCreate","title":"titleCreate"}
	
	mock.On(db.DB("content").Table("posts").Get("cc60e237-fa52-4b9c-9d72-de2ae808f535").Pluck("id", "date", "location", "author", "photo", "title", "body")).Return(
		expected, nil)

	req, err := http.NewRequest("GET", "/v1/content/cc60e237-fa52-4b9c-9d72-de2ae808f535", nil)
	if err != nil {
         t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := Handler{testEnv, Article}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}

func TestReplaceArticle(t *testing.T) {
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`
	id := `{"author":"newAuthor"}`
	post := Post{}
	json.Unmarshal([]byte(id), &post)

	var resp db.WriteResponse
	resp.Replaced = 2

	mock.On(db.DB("content").Table("posts").Get(`cc60e237-fa52-4b9c-9d72-de2ae808f535`).Replace(post)).Return(
		resp, nil)

	req, err := http.NewRequest("PUT", "/v1/content/" + articleId, bytes.NewBufferString(id))
	if err != nil {
         t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := Handler{testEnv, ReplaceArticle}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Compare bodies

	expected := `{"Errors":0,"Inserted":0,"Updated":0,"Unchanged":0,"Replaced":2,"Renamed":0,"Skipped":0,"Deleted":0,"Created":0,"DBsCreated":0,"TablesCreated":0,"Dropped":0,"DBsDropped":0,"TablesDropped":0,"GeneratedKeys":null,"FirstError":"","ConfigChanges":null,"Changes":null}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}
}

func TestUpdateArticle(t *testing.T) {
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`
	element := `author`
	newValue := `newValue`

	var resp db.WriteResponse
	resp.Replaced = 1

	mock.On(db.DB("content").Table("posts").Get(articleId).Update(`{"` + element +`": "` + newValue + `"}`)).Return(
		resp, nil)

	req, err := http.NewRequest("PATCH", "/v1/content/" + articleId + "/" + element + "/" + newValue, nil)
	if err != nil {
         t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := Handler{testEnv, UpdateArticle}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Compare bodies

	expected := `{"Errors":0,"Inserted":0,"Updated":0,"Unchanged":0,"Replaced":1,"Renamed":0,"Skipped":0,"Deleted":0,"Created":0,"DBsCreated":0,"TablesCreated":0,"Dropped":0,"DBsDropped":0,"TablesDropped":0,"GeneratedKeys":null,"FirstError":"","ConfigChanges":null,"Changes":null}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}
}

func TestDeleteArticles(t *testing.T) {
	articleId := `cc60e237-fa52-4b9c-9d72-de2ae808f535`

	mock.On(db.DB("content").Table("posts").Get(articleId).Delete()).Return(
		`{}`, nil)

	req, err := http.NewRequest("DELETE", "/v1/content/" + articleId, nil)
	if err != nil {
         t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := Handler{testEnv, DeleteArticle}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}
	
	expected :=`{}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
	}

	mock.AssertExpectations(t)
}