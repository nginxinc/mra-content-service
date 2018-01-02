package main

//
// routes.go
// ContentService
//
// Copyright © 2017 NGINX Inc. All rights reserved.
//

// Route object that stores the value of each HTTP endpoint
// Name: name of the route
// Method: HTTP request method
// Pattern: URI of endpoint
// Function: function to be associated with endpoint
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Function HandlerFunc
}

type Routes []Route

// Array of Route objects, each associated with a unique HTTP endpoint
var routes = Routes {
// Handler listening for GET at "/" URI
// Returns specified string
// @return: string "Welcome to the content service!"
	Route {
		"Welcome",
		"GET",
		"/",
		Welcome,
	},
// Handler listening for GET at "/v1/content" URI
// Get array of all articles in database
// @return: array of JSON objects, each element with single post information
	Route {
		"GetAllArticles",
		"GET",
		"/v1/content",
		GetAllArticles,
	},
// Handler listening for GET at "/v1/content/{articleID}" URI
// Get specified article
// Parameters: articleID - specifies which article to get
// @return: single JSON object with post information
	Route {
		"GetArticle",
		"GET",
		"/v1/content/{articleId}",
		GetArticle,
	},
// Handler listening for POST at "/v1/content" URI
// Creates new article based on JSON object in POST
// @POST: new post object
// @return: ID of post within database
	Route {
		"NewArticle",
		"POST",
		"/v1/content",
		NewArticle,
	},
// Handler listening for PUT at "/v1/content/{articleID}" URI
// Updates elements within specified post
// Parameters: articleID - specifies which article to update
// @POST: new post object
// @return: JSON object that specied what information was changed
	Route {
		"ReplaceArticle",
		"PUT",
		"/v1/content/{articleId}",
		ReplaceArticle,
	},
// Handler listening for PATCH at "/v1/content/{articleId}/{element}/{newValue}" URI
// Updates single element within specified post
// Parameters: articleID - specifies which article to update
// 			   element - element within post to update
// 			   newValue - new value of element
// @return: JSON object that specied what information was changed
	Route {
		"UpdateArticle",
		"PATCH",
		"/v1/content/{articleId}/{element}/{newValue}",
		UpdateArticle,
	},
// Handler listening for DELETE at "/v1/content/{articleId}" URI
// Delete specified post
// Parameters: articleID - specifies which article to delete
// @return: empty JSON object
	Route {
		"DeleteArticle",
		"DELETE",
		"/v1/content/{articleId}",
		DeleteArticle,
	},
}
