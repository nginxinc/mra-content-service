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
	"strconv"
	"crypto/tls"
)

func SetAlbumPublic(id int, public bool, r *http.Request) error {

	albumManagerHost := os.Getenv("ALBUM_MANAGER_HOST")
	albumsPath := os.Getenv("ALBUMS_PATH")
	url := albumManagerHost + albumsPath + strconv.Itoa(id) + "/public/" + strconv.FormatBool(public)
	userId := r.Header.Get("auth-id")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	albumManagerClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
		Transport: tr,
	}

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "mra-content-service")
	req.Header.Set("auth-id", userId)

	_, getErr := albumManagerClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	return nil
}
