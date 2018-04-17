package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/antchfx/htmlquery"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
)

func getNextId(db *sql.DB) int {
	nextId := 1
	rows, err := db.Query("select id from bookmarks order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		if nextId < id {
			break
		}
		nextId = nextId + 1
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return nextId
}

func addBookmark(db *sql.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		url := c.Args().First()

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("error: \"%s\" is not a proper url\n", url)
			os.Exit(1)
		}
		defer resp.Body.Close()

		html, err := htmlquery.Parse(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		title := htmlquery.InnerText(htmlquery.FindOne(html, "//title"))

		id := getNextId(db)

		uuid := uuid.Must(uuid.NewV4())

		_, err = db.Exec("insert into bookmarks(url, id, title, uuid) values(?, ?, ?, ?)", url, id, title, uuid)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("added bookmark: %s\n", title)
		return nil
	}
}

func listArchivedBookmarks(db *sql.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		rows, err := db.Query("select title, url from archive")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var url string
			var title string
			err = rows.Scan(&title, &url)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s %s\n", url, title)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
}

func listBookmarks(db *sql.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		rows, err := db.Query("select id, url, title from bookmarks")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var url string
			var id int
			var title string
			err = rows.Scan(&id, &url, &title)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%d %s %s\n", id, url, title)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
}

func archiveBookmark(db *sql.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		id := c.Args().First()
		var url string
		var uuid string
		var title string
		err := db.QueryRow("insert into archive select url, uuid, title from bookmarks where id=?", id).Scan(&url, &uuid, &title)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("delete from bookmarks where id=?", id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("archived bookmark: %s\n", title)
		return nil
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "bark"
	app.Usage = "bookmark things like you mean it"
	app.EnableBashCompletion = true
	app.Version = "0.0.1"

	barkPath := filepath.Join(os.Getenv("HOME"), ".bark")
	err := os.MkdirAll(barkPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	databaseFile := filepath.Join(barkPath, "database")
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists bookmarks (
		url text not null unique,
		id integer not null unique,
		title text,
		uuid text not null unique primary key
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
	create table if not exists archive (
		url text not null unique,
		title text,
		uuid text not null unique primary key
	);
	`

	_, err = db.Exec(sqlStmt2)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a bookmark",
			Action:  addBookmark(db),
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list bookmarks",
			Action:  listBookmarks(db),
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "remove"},
			Usage:   "archive a bookmark",
			Action:  archiveBookmark(db),
		},
		{
			Name:   "archive",
			Usage:  "list archived bookmarks",
			Action: listArchivedBookmarks(db),
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
