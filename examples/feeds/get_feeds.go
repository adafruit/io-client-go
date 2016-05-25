package main

// Run with:
//   go run get_feeds.go -key "MY AIO KEY"

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

	flag.StringVar(&useURL, "url", "http://localhost:3002", "Adafruit IO URL")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")

	flag.Parse()

	client := adafruitio.NewClient(key)
	client.BaseURL, _ = url.Parse(useURL)

	// Get the list of all available feeds
	feeds, _, err := client.Feed.All()
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	for _, feed := range feeds {
		sfeed, _ := json.Marshal(feed)
		fmt.Println(string(sfeed))
	}
}
