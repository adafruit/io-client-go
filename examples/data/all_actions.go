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
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&useURL, "url", "", "Adafruit IO URL")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")
	flag.StringVar(&feedName, "feed", "", "the key of the feed to send to")
	flag.StringVar(&value, "value", "", "the value to send")

	if useURL == "" {
		// no arg given, try ENV
		useURL = os.Getenv("ADAFRUIT_IO_URL")
	}

	if key == "" {
		key = os.Getenv("ADAFRUIT_IO_KEY")
	}

	if value == "" {
		value = rval()
	}

	if feedName == "" {
		// generate feed name
		feedName = fmt.Sprintf("beta-test-%v", fmt.Sprintf("%06d", rand.Int())[0:6])
	}

	flag.Parse()
}

func rval() string {
	return fmt.Sprintf("%f", rand.Float32()*100.0)
}

func render(label string, d *aio.DataPoint) {
	dbytes, _ := json.MarshalIndent(d, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(dbytes))
}

func title(label string) {
	fmt.Printf("\n\n%v\n\n", label)
}

func main() {
	prepare()

	client := aio.NewClient(key)
	client.BaseURL, _ = url.Parse(useURL)

	feed, _, ferr := client.Feed.Get(feedName)
	if ferr != nil {
		fmt.Printf("unable to load Feed with key %v, creating placeholder\n", feedName)
		feed = &aio.Feed{Key: feedName}
	}

	// create a data point on an existing Feed, create Feed if needed
	client.SetFeed(feed)
	val := &aio.DataPoint{Value: json.Number(value)}

	title("Create and Check")

	dp, _, err := client.Data.Send(val)
	if err != nil {
		fmt.Println("unable to send data")
		panic(err)
	}
	render("new point", dp)

	ndp, _, err := client.Data.Get(dp.ID)
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("found point", ndp)

	// update point
	client.Data.Update(dp.ID, &aio.DataPoint{Value: json.Number(rval())})

	// reload
	ndp, _, err = client.Data.Get(dp.ID)
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("updated point", ndp)
	time.Sleep(2 * time.Second)

	// Generate some more Data to fill out the stream
	for i := 0; i < 4; i += 1 {
		client.Data.Create(&aio.DataPoint{Value: json.Number(rval())})
	}

	// Display all Data in the stream
	title("All Data")
	dts, _, err := client.Data.All()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.ID), data)
	}
	time.Sleep(2 * time.Second)

	// stream commands: Last, Prev, and Next
	title("Queue related commands")

	ndp, _, err = client.Data.Last()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("last point", ndp)

	ndp, _, err = client.Data.Prev()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("prev point", ndp)

	ndp, _, err = client.Data.Next()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("next point", ndp)

	// delete
	title("Delete")
	_, derr := client.Data.Delete(ndp.ID)
	if derr == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("failed to delete!")
	}
	time.Sleep(1 * time.Second)

	title("All Data (updated)")
	dts, _, err = client.Data.All()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.ID), data)
	}
	time.Sleep(2 * time.Second)
}
