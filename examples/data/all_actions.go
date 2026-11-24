// Demo showing Data listing, creation, updating, and deletion.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	adafruitio "github.com/adafruit/io-client-go/v3"
)

var (
	useURL   string
	username string
	key      string
	feedName string
	value    string
)

func prepare() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&useURL, "url", "", "Adafruit IO URL")
	flag.StringVar(&username, "user", "", "your Adafruit IO user name")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")
	flag.StringVar(&feedName, "feed", "", "the key of the feed to manipulate")
	flag.StringVar(&value, "value", rval(), "the value to send")

	if useURL == "" {
		// no arg given, try ENV
		useURL = os.Getenv("ADAFRUIT_IO_URL")
	}

	if key == "" {
		key = os.Getenv("ADAFRUIT_IO_KEY")
	}

	if username == "" {
		username = os.Getenv("ADAFRUIT_IO_USERNAME")
	}

	flag.Parse()

	if feedName == "" {
		panic("A feed name must be specified")
	}

}

func rval() string {
	return fmt.Sprintf("%f", rand.Float32()*100.0)
}

func render(label string, d *adafruitio.Data) {
	dbytes, _ := json.MarshalIndent(d, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(dbytes))
}

func title(label string) {
	fmt.Printf("\n\n%v\n\n", label)
}

func main() {
	ctx := context.Background()

	prepare()

	client := adafruitio.NewClient(username, key)
	if useURL != "" {
		client.SetBaseURL(useURL)
	}
	feed, _, ferr := client.Feed.Get(ctx, feedName)
	if ferr != nil {
		panic(ferr)
	}

	// create a data point on an existing Feed
	val := &adafruitio.Data{Value: adafruitio.Ptr(value), FeedKey: adafruitio.Ptr(feedName)}

	title("Create and Check")

	dp, _, err := client.Data.Create(ctx, feed.GetKey(), val)
	if err != nil {
		fmt.Println("unable to create data")
		panic(err)
	}
	render("new point", dp)

	ndp, _, err := client.Data.Get(ctx, feed.GetKey(), dp.GetID())
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("found point", ndp)

	// update point
	client.Data.Update(ctx, feed.GetKey(), dp.GetID(), &adafruitio.Data{Value: adafruitio.Ptr(rval())})

	// reload
	ndp, _, err = client.Data.Get(ctx, feed.GetKey(), dp.GetID())
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("updated point", ndp)
	time.Sleep(2 * time.Second)

	// Generate some more Data to fill out the stream
	for i := 0; i < 4; i += 1 {
		client.Data.Create(ctx, feed.GetKey(), &adafruitio.Data{Value: adafruitio.Ptr(rval())})
	}

	// Display all Data in the stream
	title("All Data")
	dts, _, err := client.Data.All(ctx, feed.GetKey(), nil)
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.GetID()), data)
	}
	time.Sleep(2 * time.Second)

	// stream commands: Last, Prev, and Next
	title("Queue related commands")

	ndp, _, err = client.Data.Last(ctx, feed.GetKey())
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("last point", ndp)

	ndp, _, err = client.Data.Prev(ctx, feed.GetKey())
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("prev point", ndp)

	ndp, _, err = client.Data.First(ctx, feed.GetKey())
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("first point", ndp)

	ndp, _, err = client.Data.Next(ctx, feed.GetKey())
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("next point", ndp)

	// delete
	title("Delete")
	_, derr := client.Data.Delete(ctx, feed.GetKey(), ndp.GetID())
	if derr == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("failed to delete!")
	}
	time.Sleep(1 * time.Second)

	title("All Data (updated)")
	dts, _, err = client.Data.All(ctx, feed.GetKey(), nil)
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.GetID()), data)
	}
	time.Sleep(2 * time.Second)

	// Now, generate a single point and do a filtered search for it
	t := time.Now().Unix() // get current time
	time.Sleep(2 * time.Second)
	client.Data.Create(ctx, feed.GetKey(), &adafruitio.Data{Value: adafruitio.Ptr(rval())}) // create point 2 seconds later

	title(fmt.Sprintf("Filtered Data, since %v", t))
	dts, _, err = client.Data.All(ctx, feed.GetKey(), &adafruitio.DataFilter{StartTime: fmt.Sprintf("%d", t)})
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.GetID()), data)
	}
}
