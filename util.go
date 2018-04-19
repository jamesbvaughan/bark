package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

func getBookmark(bookmarks []bookmark, iString string) bookmark {
	i, err := strconv.Atoi(iString)
	if err != nil {
		log.Fatal(err)
	}
	return bookmarks[i-1]
}

func getPageTitle(url string) (title string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: \"%s\" is not a proper url\n", url)
		return
	}
	defer resp.Body.Close()

	html, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return
	}

	titleElement := htmlquery.FindOne(html, "//title")
	title = strings.TrimSpace(htmlquery.InnerText(titleElement))

	return
}
