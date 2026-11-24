package adafruitio_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	adafruitio "github.com/adafruit/io-client-go/v3"
)

var (
	key   string
	feeds []*adafruitio.Feed
)

func Example() {
	ctx := context.Background()

	// Load ADAFRUIT_IO_KEY from environment
	client := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_USERNAME"), os.Getenv("ADAFRUIT_IO_KEY"))

	// set custom API URL
	client.SetBaseURL("http://localhost:3002")

	// Get the list of all available feeds
	feeds, _, err := client.Feed.All(ctx)
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
