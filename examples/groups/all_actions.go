// Demo showing Group listing, creation, updating, and deletion.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	adafruitio "github.com/adafruit/io-client-go/v2"
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

func ShowAll(client *adafruitio.Client) {
	title("All")
	groups, _, err := client.Group.All()
	if err != nil {
		panic(err)
	}
	for _, g := range groups {
		render(g.Name, g)
	}
}

func pause() {
	time.Sleep(2 * time.Second)
}

func main() {
	prepare()

	client := adafruitio.NewClient(username, key)
	if useURL != "" {
		client.SetBaseURL(useURL)
	}

	ShowAll(client)
	pause()

	// CREATE
	name := fmt.Sprintf("a_new_group_%d", rand.Int())
	fmt.Printf("CREATING %v\n", name)
	g, resp, err := client.Group.Create(&adafruitio.Group{Name: name})
	if err != nil {
		// resp.Debug()
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
	newg, _, err := client.Group.Get(g.Key)
	if err != nil {
		panic(err)
	}
	render("new group", newg)
	pause()

	// UPDATE (only Name and Description can be modified)
	g.Name = fmt.Sprintf("name_changed_to_%d", rand.Int())
	g.Description = "Now this group has a description."

	fmt.Printf("changing name to %v\n", g.Name)
	newg, _, err = client.Group.Update(g.Key, g)
	if err != nil {
		panic(err)
	}
	render("updated group", newg)
	pause()

	// DELETE
	time.Sleep(2 * time.Second)
	title("deleting group")
	_, err = client.Group.Delete(newg.Key)
	if err == nil {
		fmt.Println("ok")
	}
	pause()

	// SHOW ALL
	ShowAll(client)
}
