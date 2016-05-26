package main

// Run with:
//   go run send_data.go -key "MY AIO KEY" -feed "temp" -value "98.6"

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"time"

	aio "github.com/adafruit/io-client-go"
)

var (
	useURL   string
	key      string
	feedName string
	value    string
)

func prepare() {
	flag.StringVar(&useURL, "url", "", "Adafruit IO URL")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")
	flag.StringVar(&feedName, "feed", "beta-test", "the key of the feed to send to")
	flag.StringVar(&value, "value", "", "the value to send")

	if useURL == "" {
		// no arg given, try ENV
		useURL = os.Getenv("ADAFRUIT_IO_URL")
	}

	if key == "" {
		key = os.Getenv("ADAFRUIT_IO_KEY")
	}

	if value == "" {
		rand.Seed(time.Now().UnixNano())
		value = fmt.Sprintf("%f", rand.Float32()*100.0)
	}

	flag.Parse()
}

func render(label string, d *aio.DataPoint) {
	dbytes, _ := json.MarshalIndent(d, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(dbytes))
}

func main() {
	prepare()

	client := aio.NewClient(key)
	client.BaseURL, _ = url.Parse(useURL)

	feed, _, ferr := client.Feed.Get(feedName)
	if ferr != nil {
		fmt.Println("unable to load Feed, creating placeholder")
		feed = &aio.Feed{Key: feedName}
	}

	// create a data point on an existing Feed, create Feed if needed
	client.SetFeed(feed)
	val := &aio.DataPoint{Value: json.Number(value)}

	dp, _, err := client.Data.Send(val)
	if err != nil {
		fmt.Println("unable to send data")
		panic(err)
	}

	render("new point", dp)

	dts, _, err := client.Data.All()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("\n\n----==== ALL DATA POINTS ====----\n")
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.ID), data)
	}
}
