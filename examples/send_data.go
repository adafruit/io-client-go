package main

// Run with:
//   go run send_data.go -key "MY AIO KEY"

import (
	// provides adafruitio
	"encoding/json"
	"flag"
	"fmt"

	"github.com/adafruit/io-client-go"
)

func main() {
	var key string
	var feed string
	var value string

	flag.StringVar(&key, "key", "", "your Adafruit IO key")
	flag.StringVar(&feed, "feed", "beta-test", "the feed to send to")
	flag.StringVar(&value, "value", "42", "the value to send")

	flag.Parse()

	client := adafruitio.NewClient(key)
	// client.BaseURL, _ = url.Parse("http://localhost:3002/")

	// create a data point on an existing Feed
	client.SetFeed(feed)
	val := json.Number(value)

	dp, _, err := client.Data.Send(&adafruitio.DataPoint{Value: &val})
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	fmt.Println("generated datapoint:", dp)
}
