// Demo showing Data listing, creation, updating, and deletion.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/adafruit/io-client-go/v2/pkg/adafruitio"
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
	prepare()

	client := adafruitio.NewClient(username, key)
	if useURL != "" {
		client.SetBaseURL(useURL)
	}
	feed, _, ferr := client.Feed.Get(feedName)
	if ferr != nil {
		panic(ferr)
	}

	// create a data point on an existing Feed
	client.SetFeed(feed)
	val := &adafruitio.Data{Value: value, FeedKey: feedName}

	title("Create and Check")

	dp, _, err := client.Data.Create(val)
	if err != nil {
		fmt.Println("unable to create data")
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
	client.Data.Update(dp.ID, &adafruitio.Data{Value: rval()})

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
		client.Data.Create(&adafruitio.Data{Value: rval()})
	}

	// Display all Data in the stream
	title("All Data")
	dts, _, err := client.Data.All(nil)
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

	ndp, _, err = client.Data.First()
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	render("first point", ndp)

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
	dts, _, err = client.Data.All(nil)
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.ID), data)
	}
	time.Sleep(2 * time.Second)

	// Now, generate a single point and do a filtered search for it
	t := time.Now().Unix() // get current time
	time.Sleep(2 * time.Second)
	client.Data.Create(&adafruitio.Data{Value: rval()}) // create point 2 seconds later

	title(fmt.Sprintf("Filtered Data, since %v", t))
	dts, _, err = client.Data.All(&adafruitio.DataFilter{StartTime: fmt.Sprintf("%d", t)})
	if err != nil {
		fmt.Println("unable to retrieve data")
		panic(err)
	}
	for _, data := range dts {
		render(fmt.Sprintf("ID <%v>", data.ID), data)
	}
}
