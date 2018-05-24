#!/bin/sh

echo "Running Tests"

ALBUM_MANAGER_HOST=http://localhost:8080 ALBUMS_PATH=/albums go test -v ''
