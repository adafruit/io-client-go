// Demo showing feed listing, creation, updating, and deletion.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	adafruitio "github.com/adafruit/io-client-go/v3"
)

var (
	useURL   string
	username string
	key      string
)

func prepare() {
	flag.StringVar(&useURL, "url", "", "Adafruit IO URL")
	flag.StringVar(&username, "user", "", "your Adafruit IO user name")
	flag.StringVar(&key, "key", "", "your Adafruit IO key")

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
}

func render(label string, f *adafruitio.Feed) {
	sfeed, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(sfeed))
}

func title(label string) {
	fmt.Printf("\n\n%v\n\n", label)
}

func ShowAll(ctx context.Context, client *adafruitio.Client) {
	// Get the list of all available feeds
	feeds, _, err := client.Feed.All(ctx)
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	for _, feed := range feeds {
		render(feed.GetName(), feed)
	}
}

func main() {
	ctx := context.Background()

	prepare()

	client := adafruitio.NewClient(username, key)
	if useURL != "" {
		client.SetBaseURL(useURL)
	}

	title("All")

	ShowAll(ctx, client)
	time.Sleep(1 * time.Second)

	title("Create")

	newFeed := &adafruitio.Feed{Name: adafruitio.Ptr("my-new-feed"), Description: adafruitio.Ptr("an example of creating feeds")}
	client.Feed.Create(ctx, newFeed)
	render("NEW FEED", newFeed)
	time.Sleep(1 * time.Second)

	if newFeed.GetID() == 0 {
		panic(fmt.Errorf("could not create feed"))
	}

	title("All")

	ShowAll(ctx, client)
	time.Sleep(1 * time.Second)

	title("Update")

	updatedFeed, _, _ := client.Feed.Update(ctx, newFeed.GetKey(), &adafruitio.Feed{Name: adafruitio.Ptr("renamed-new-feed")})
	render("UPDATED FEED", updatedFeed)
	time.Sleep(1 * time.Second)

	title("All")

	ShowAll(ctx, client)
	time.Sleep(1 * time.Second)

	title("Just the new one")

	feed, _, _ := client.Feed.Get(ctx, updatedFeed.GetKey())
	render("SHOW", feed)
	time.Sleep(1 * time.Second)

	title("Delete")

	_, err := client.Feed.Delete(ctx, newFeed.GetKey())
	if err == nil {
		fmt.Println("ok")
	}
	time.Sleep(1 * time.Second)

	title("All")

	ShowAll(ctx, client)
	time.Sleep(1 * time.Second)
}
