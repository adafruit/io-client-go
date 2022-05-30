// A simple web service that displays data from your Adafruit IO feeds.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	adafruitio "github.com/adafruit/io-client-go/v2"
)

var (
	aioURL      string
	aioUsername string
	aioKey      string
	feedMatcher = regexp.MustCompile(`/feed/([a-z0-9-]+)`)
	head        = `
	<!doctype html><head>
		<style>
			body, div, h1 { margin: 0; padding: 0; border: 0; font-size: 100%; font: inherit; vertical-align: baseline; }
			table { border-collapse: collapse; border-spacing: 0; }
			td, th { padding: 8px; }
			th { border-bottom: 1px solid #888; }
			.wrap { width: 640px; margin: 0 auto; }
		</style>
	</head><body><div class='wrap'>`
	tail = `</div></body>`
)

// Get command line flags, fallback to environment variables
func prepare() {
	flag.StringVar(&aioURL, "url", "", "Adafruit IO URL")
	flag.StringVar(&aioUsername, "user", "", "your Adafruit IO user name")
	flag.StringVar(&aioKey, "key", "", "your Adafruit IO key")

	if aioURL == "" {
		// no arg given, try ENV
		aioURL = os.Getenv("ADAFRUIT_IO_URL")
	}

	if aioKey == "" {
		aioKey = os.Getenv("ADAFRUIT_IO_KEY")
	}

	if aioUsername == "" {
		aioUsername = os.Getenv("ADAFRUIT_IO_USERNAME")
	}

	flag.Parse()
}

func main() {
	prepare()

	// setup AIO client
	client := adafruitio.NewClient(aioUsername, aioKey)

	// setup server
	mux := http.NewServeMux()

	// GET /
	// get feeds list
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, head)
			defer fmt.Fprintf(w, tail)

			fmt.Fprint(w, "<h1>Feeds</h1>")

			feeds, _, err := client.Feed.All()
			if err != nil {
				fmt.Fprintf(w, "ERROR finding feed. %v", err.Error())
				return
			}

			fmt.Fprint(w, "<ul>")
			for _, feed := range feeds {
				fmt.Fprintf(w,
					`<li><a href='/feed/%v'>%v</a></li>`,
					feed.Key,
					feed.Name,
				)
			}
			fmt.Fprint(w, "</ul>")
		},
	)

	// GET /feed/{feedname}
	// get data for a given feed
	mux.HandleFunc("/feed/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, head)
			defer fmt.Fprintf(w, tail)

			var parts []string
			var feedName string

			parts = feedMatcher.FindStringSubmatch(r.URL.Path)
			if len(parts) > 0 {
				feedName = parts[1]
				fmt.Fprintf(w, "<h1>Feed: %v</h1><a href='/'>back</a>", feedName)
			} else {
				fmt.Fprint(w, "ERROR: unable to find feed")
				return
			}

			feed, _, err := client.Feed.Get(feedName)

			if err != nil {
				fmt.Fprintf(w, "ERROR finding feed. %v", err.Error())
				return
			}

			// set client Feed
			client.SetFeed(feed)

			// get all Data for the given Feed
			data, _, err := client.Data.All(nil)
			if err != nil {
				fmt.Fprintf(w, "ERROR loading data. %v", err.Error())
				return
			}

			// render data in a table
			fmt.Fprint(w, "<table><tr><th>Created At</th><th>Value</th></tr>")
			for _, d := range data {
				fmt.Fprintf(w,
					`<tr>
						<td>%v</td>
						<td>%v</td>
					</tr>`,
					d.CreatedAt,
					d.Value,
				)
			}
			fmt.Fprint(w, "</table>")
		},
	)

	// start server
	log.Println("Listening on port 3009...")
	http.ListenAndServe(":3009", mux)
}
