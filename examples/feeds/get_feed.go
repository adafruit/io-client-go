package main

// Run with:
//   go run get_feed.go -key "MY AIO KEY" -feed "feed name, key, or ID"

import (
	// provides adafruitio

	"encoding/json"
	"flag"
	"fmt"
	"net/url"

	"github.com/adafruit/io-client-go"
)

func main() {
	var useURL string
	var key string
	var feedID string

	flag.StringVar(&useURL, "url", "http://localhost:3002", "Adafruit IO URL")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")
	flag.StringVar(&feedID, "feed", "beta-test", "the feed to look up")

	flag.Parse()

	client := adafruitio.NewClient(key)
	client.BaseURL, _ = url.Parse(useURL)

	// Get a single feed
	feed, _, err := client.Feed.Get(feedID)
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	sfeed, _ := json.Marshal(feed)
	fmt.Println(string(sfeed))
}
