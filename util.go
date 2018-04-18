package main

import (
	"log"
	"strconv"
)

func getBookmark(bookmarks []bookmark, iString string) bookmark {
	i, err := strconv.Atoi(iString)
	if err != nil {
		log.Fatal(err)
	}
	return bookmarks[i-1]
}
