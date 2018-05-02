package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PageData struct {
	Bookmarks []Bookmark
	Page      string
}

func Serve() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("web/assets/"))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bookmarks := GetBookmarks()

		data := PageData{
			Bookmarks: bookmarks,
			Page:      "bookmarks",
		}

		tmpl := template.Must(template.ParseFiles("web/templates/layout.html"))
		tmpl.Execute(w, data)
	})

	r.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
		bookmarks := GetArchivedBookmarks()

		data := PageData{
			Bookmarks: bookmarks,
			Page:      "archive",
		}

		tmpl := template.Must(template.ParseFiles("web/templates/layout.html"))
		tmpl.Execute(w, data)
	})

	r.HandleFunc("/bookmarks/{uuid}/archive", func(w http.ResponseWriter, r *http.Request) {
		bookmarks := GetBookmarks()

		vars := mux.Vars(r)
		uuid := vars["uuid"]

		for i, bookmark := range bookmarks {
			if bookmark.UUID == uuid {
				log.Printf("Archiving bookmark %d: %s\n", i, bookmark.Title)
				ArchiveBookmark(bookmark.UUID)
				break
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.HandleFunc("/bookmarks/{uuid}/delete", func(w http.ResponseWriter, r *http.Request) {
		bookmarks := GetArchivedBookmarks()

		vars := mux.Vars(r)
		uuid := vars["uuid"]

		for _, bookmark := range bookmarks {
			if bookmark.UUID == uuid {
				log.Printf("Deleting bookmark: %s\n", bookmark.Title)
				DeleteBookmark(bookmark.UUID)
				break
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	fmt.Println("listening on http://localhost:3030")
	http.ListenAndServe(":3030", r)
}
