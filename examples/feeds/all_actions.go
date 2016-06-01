// Demo showing feed listing, creation, updating, and deletion.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	aio "github.com/adafruit/io-client-go"
)

var (
	useURL string
	key    string
)

func prepare() {
	flag.StringVar(&useURL, "url", "", "Adafruit IO URL")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")

	if useURL == "" {
		// no arg given, try ENV
		useURL = os.Getenv("ADAFRUIT_IO_URL")
	}

	if key == "" {
		key = os.Getenv("ADAFRUIT_IO_KEY")
	}

	flag.Parse()
}

func render(label string, f *aio.Feed) {
	sfeed, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(sfeed))
}

func title(label string) {
	fmt.Printf("\n\n%v\n\n", label)
}

func ShowAll(client *aio.Client) {
	// Get the list of all available feeds
	feeds, _, err := client.Feed.All()
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	for _, feed := range feeds {
		render(feed.Name, feed)
	}
}

func main() {
	prepare()

	client := aio.NewClient(key)
	client.BaseURL, _ = url.Parse(useURL)

	title("All")

	ShowAll(client)
	time.Sleep(1 * time.Second)

	title("Create")

	newFeed := &aio.Feed{Name: "my-new-feed", Description: "an example of creating feeds"}
	client.Feed.Create(newFeed)
	render("NEW FEED", newFeed)
	time.Sleep(1 * time.Second)

	if newFeed.ID == 0 {
		panic(fmt.Errorf("could not create feed"))
	}

	title("All")

	ShowAll(client)
	time.Sleep(1 * time.Second)

	title("Update")

	updatedFeed, _, _ := client.Feed.Update(newFeed.ID, &aio.Feed{Name: "renamed-new-feed"})
	render("UPDATED FEED", updatedFeed)
	time.Sleep(1 * time.Second)

	title("All")

	ShowAll(client)
	time.Sleep(1 * time.Second)

	title("Just the new one")

	feed, _, _ := client.Feed.Get(updatedFeed.ID)
	render("SHOW", feed)
	time.Sleep(1 * time.Second)

	title("Delete")

	_, err := client.Feed.Delete(newFeed.ID)
	if err == nil {
		fmt.Println("ok")
	}
	time.Sleep(1 * time.Second)

	title("All")

	ShowAll(client)
	time.Sleep(1 * time.Second)
}
