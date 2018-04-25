package bark

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/antchfx/htmlquery"
)

func GetBookmark(bookmarks []Bookmark, indexString string) Bookmark {
	i, err := strconv.Atoi(indexString)
	if err != nil {
		log.Fatal(err)
	}
	return bookmarks[i-1]
}

func GetPageTitle(inputURL string) (title string, err error) {
	resp, err := http.Get(inputURL)
	if err != nil {
		fmt.Printf("error: \"%s\" is not a proper url\n", inputURL)
		return
	}
	defer resp.Body.Close()

	html, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return
	}

	titleElement := htmlquery.FindOne(html, "//title")
	if titleElement == nil {
		u, err := url.Parse(inputURL)
		if err != nil {
			log.Fatal(err)
		}

		title = u.Path
	} else {
		title = strings.TrimSpace(htmlquery.InnerText(titleElement))
	}

	return
}

func GetHostFromURL(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatal(err)
	}

	return u.Hostname()
}

func PrintBookmarkTable(bookmarks []Bookmark, printURLs bool, printIDs bool) {
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
		fmt.Fprintf(w, "%s", bookmark.Title)
		if printURLs {
			fmt.Fprintf(w, "\t%s", bookmark.URL)
		}
		fmt.Fprintln(w)
	}
}
