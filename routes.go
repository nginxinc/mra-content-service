package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes {
	Route {
		"Welcome",
		"GET",
		"/",
		Welcome,
	},
	Route {
		"Articles",
		"GET",
		"/v1/content",
		Articles,
	},
	Route {
		"Article",
		"GET",
		"/v1/content/{articleId}",
		Article,
	},
	Route {
		"NewArticle",
		"POST",
		"/v1/content",
		NewArticle,
	},
	Route {
		"ReplaceArticle",
		"PUT",
		"/v1/content/{aticleId}",
		ReplaceArticle,
	},
	Route {
		"UpdateArticle",
		"PUT",
		"/v1/content/{articleId}/{element}/{newValue}",
		UpdateArticle,
	},
	Route {
		"DeleteArticle",
		"DELETE",
		"/v1/content/{articleId}",
		DeleteArticle,
	},
}
