package main

//
// album_manager.go
// ContentService
//
// Copyright © 2018 NGINX Inc. All rights reserved.
//

import (
	"log"
	"net/http"
	"time"
	"os"
)

type album struct {
	Number int `json:"number"`
}

func SetAlbumPublic(id string, r *http.Request) error {

	albumManagerHost := os.Getenv("ALBUM_MANAGER_HOST")
	albumsPath := os.Getenv("ALBUMS_PATH")
	url := albumManagerHost + albumsPath + id + "/public/true"
	userId := r.Header.Get("auth-id")

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "mra-content-service")
	req.Header.Set("auth-id", userId)

	_, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return getErr
}

func SetAlbumPrivate(id string, r *http.Request) error {

	albumManagerHost := os.Getenv("ALBUM_MANAGER_HOST")
	albumsPath := os.Getenv("ALBUMS_PATH")
	url := albumManagerHost + albumsPath + id + "/public/false"
	userId := r.Header.Get("auth-id")

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "mra-content-service")
	req.Header.Set("auth-id", userId)

	_, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return getErr
}