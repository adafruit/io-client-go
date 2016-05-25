package adafruitio_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/feeds",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `[{"id":1, "name":"feed"}]`)
		},
	)

	assert := assert.New(t)

	feeds, response, err := client.Feed.All()

	assert.Nil(err)
	assert.NotNil(feeds)
	assert.NotNil(response)

	feed := feeds[0]

	assert.Equal(1, feed.ID)
	assert.Equal("feed", feed.Name)
}

func TestFeedGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/feeds/test",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":1, "name":"test"}`)
		},
	)

	assert := assert.New(t)

	feed, response, err := client.Feed.Get("test")

	assert.Nil(err)
	assert.NotNil(feed)
	assert.NotNil(response)

	assert.Equal(1, feed.ID)
	assert.Equal("test", feed.Name)
}
