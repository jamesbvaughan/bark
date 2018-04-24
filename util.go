package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/antchfx/htmlquery"
)

func getBookmark(bookmarks []bookmark, indexString string) bookmark {
	i, err := strconv.Atoi(indexString)
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

func printBookmarkTable(bookmarks []bookmark, printURLs bool, printIDs bool) {
	w := tabwriter.NewWriter(os.Stdout, 2, 1, 1, ' ', 0)
	defer w.Flush()

	if printIDs {
		fmt.Fprint(w, "ID\t")
	}
	fmt.Fprint(w, "Title")
	if printURLs {
		fmt.Fprint(w, "\tURL")
	}
	fmt.Fprintln(w)

	if printIDs {
		fmt.Fprint(w, "--\t")
	}
	fmt.Fprint(w, "-----")
	if printURLs {
		fmt.Fprint(w, "\t---")
	}
	fmt.Fprintln(w)

	for i, bookmark := range bookmarks {
		if printIDs {
			fmt.Fprintf(w, "%d\t", i+1)
		}
		fmt.Fprintf(w, "%s", bookmark.title)
		if printURLs {
			fmt.Fprintf(w, "\t%s", bookmark.url)
		}
		fmt.Fprintln(w)
	}
}
