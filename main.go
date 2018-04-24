package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "bark"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "James Vaughan",
			Email: "james@jamesbvaughan.com",
		},
	}
	app.Copyright = "(c) 2018 James Vaughan"
	app.HelpName = "bark"
	app.Usage = "bookmark things like you mean it"
	app.EnableBashCompletion = true

	err := initializeDatabase()
	if err != nil {
		log.Fatal(err)
	}

	bookmarks, err := getBookmarks()
	if err != nil {
		log.Fatal(err)
	}

	archivedBookmarks, err := getArchivedBookmarks()
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a bookmark",
			Action: func(c *cli.Context) (err error) {
				url := c.Args().First()
				title, err := addBookmark(url)
				if err != nil {
					return
				}

				fmt.Printf("added bookmark %d: %s\n", len(bookmarks)+1, title)
				return
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list bookmarks",
			UseShortOptionHandling: true,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "urls, u",
					Usage: "print bookmark URLs",
				},
				cli.BoolFlag{
					Name:  "archive, a",
					Usage: "print archived bookmarks",
				},
			},
			Action: func(c *cli.Context) (err error) {
				bookmarksToPrint := bookmarks
				if c.Bool("archive") {
					bookmarksToPrint = archivedBookmarks
				}
				printBookmarkTable(bookmarksToPrint, c.Bool("urls"), !c.Bool("archive"))
				return
			},
		},
		{
			Name:    "open",
			Aliases: []string{"o"},
			Usage:   "open a bookmark",
			Action: func(c *cli.Context) (err error) {
				bookmark := getBookmark(bookmarks, c.Args().First())
				fmt.Printf("opening \"%s\"...\n", bookmark.title)
				err = open.Run(bookmark.url)
				return
			},
		},
		{
			Name:  "archive",
			Usage: "archive a bookmark",
			Action: func(c *cli.Context) (err error) {
				bookmark := getBookmark(bookmarks, c.Args().First())
				err = archiveBookmark(bookmark)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("archived bookmark: \"%s\"\n", bookmark.title)
				return
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "rm"},
			Usage:   "permanently delete a bookmark",
			Action: func(c *cli.Context) (err error) {
				bookmark := getBookmark(bookmarks, c.Args().First())
				err = deleteBookmark(bookmark)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("deleted bookmark: \"%s\"\n", bookmark.title)
				return
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
