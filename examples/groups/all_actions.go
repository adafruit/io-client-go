// Demo showing Group listing, creation, updating, and deletion.
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
)

func prepare() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&useURL, "url", "", "Adafruit IO URL. Defaults to https://io.adafruit.com.")
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

func render(label string, f *adafruitio.Group) {
	sfeed, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(sfeed))
}

func title(label string) {
	fmt.Printf("\n\n%v\n\n", label)
}

func ShowAll(ctx context.Context, client *adafruitio.Client) {
	title("All")
	groups, _, err := client.Group.All(ctx)
	if err != nil {
		panic(err)
	}
	for _, g := range groups {
		render(g.GetName(), g)
	}
}

func pause() {
	time.Sleep(2 * time.Second)
}

func main() {
	ctx := context.Background()

	prepare()

	client := adafruitio.NewClient(username, key)
	if useURL != "" {
		client.SetBaseURL(useURL)
	}

	ShowAll(ctx, client)
	pause()

	// CREATE
	name := fmt.Sprintf("a_new_group_%d", rand.Int())
	fmt.Printf("CREATING %v\n", name)
	g, resp, err := client.Group.Create(ctx, &adafruitio.Group{Name: adafruitio.Ptr(name)})
	if err != nil {
		fmt.Printf("failed to create group")
		panic(err)
	} else if resp.StatusCode > 299 {
		fmt.Printf("Unexpected status: %v", resp.Status)
		panic(fmt.Errorf("failed to create group"))
	} else {
		fmt.Println("ok")
	}
	pause()

	// GET
	newg, _, err := client.Group.Get(ctx, g.GetKey())
	if err != nil {
		panic(err)
	}
	render("new group", newg)
	pause()

	// UPDATE (only Name and Description can be modified)
	g.Name = adafruitio.Ptr(fmt.Sprintf("name_changed_to_%d", rand.Int()))
	g.Description = adafruitio.Ptr("Now this group has a description.")

	fmt.Printf("changing name to %v\n", g.GetName())
	newg, _, err = client.Group.Update(ctx, g.GetKey(), g)
	if err != nil {
		panic(err)
	}
	render("updated group", newg)
	pause()

	// DELETE
	time.Sleep(2 * time.Second)
	title("deleting group")
	_, err = client.Group.Delete(ctx, newg.GetKey())
	if err == nil {
		fmt.Println("ok")
	}
	pause()

	// SHOW ALL
	ShowAll(ctx, client)
}
