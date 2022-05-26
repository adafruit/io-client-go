package adafruitio_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/adafruit/io-client-go"
)

var (
	key   string
	feeds []*adafruitio.Feed
)

func Example() {
	// Load ADAFRUIT_IO_KEY from environment
	client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_KEY"))

	// set custom API URL
	client.BaseURL, _ = url.Parse("http://localhost:3002")

	// Get the list of all available feeds
	feeds, _, err := client.Feed.All()
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	// View the resulting feed list
	for _, feed := range feeds {
		jsonBytes, _ := json.MarshalIndent(feed, "", "  ")
		fmt.Printf("[%v]\n", feed.Name)
		fmt.Println(string(jsonBytes))
	}
}
