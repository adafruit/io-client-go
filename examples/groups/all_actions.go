package main

// Run with:
//   go run get_feed.go -key "MY AIO KEY" -feed "feed name, key, or ID"

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
	useURL string
	key    string
)

func prepare() {
	rand.Seed(time.Now().UnixNano())

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

func render(label string, f *aio.Group) {
	sfeed, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(sfeed))
}

func title(label string) {
	fmt.Printf("\n\n%v\n\n", label)
}

func ShowAll(client *aio.Client) {
	title("All")
	groups, _, err := client.Group.All()
	if err != nil {
		panic(err)
	}
	for _, g := range groups {
		render(g.Name, g)
	}
}

func main() {
	prepare()

	client := aio.NewClient(key)
	client.BaseURL, _ = url.Parse(useURL)

	ShowAll(client)

	var g *aio.Group
	var name string

	name = fmt.Sprintf("a_new_group_%d", rand.Int())
	fmt.Printf("CREATING %v\n", name)
	g, resp, err := client.Group.Create(&aio.Group{Name: name})
	if err != nil {
		// resp.Debug()
		fmt.Printf("failed to create group")
		panic(err)
	} else if resp.StatusCode > 299 {
		fmt.Printf("Unexpected status: %v", resp.Status)
		panic(fmt.Errorf("failed to create group"))
	}
	render("new group", g)
}
