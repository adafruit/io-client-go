package adafruitio

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds"),
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

	mux.HandleFunc(serverPattern("feeds/test"),
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

func TestFeedCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":1, "name":"test"}`)
		},
	)

	assert := assert.New(t)

	nfeed := &Feed{Name: "test"}

	feed, response, err := client.Feed.Create(nfeed)

	assert.Nil(err)
	assert.NotNil(feed)
	assert.NotNil(response)

	assert.Equal(1, feed.ID)
	assert.Equal("test", feed.Name)
}

func TestFeedUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PATCH")
			fmt.Fprint(w, `{"id":1, "name":"test-2"}`)
		},
	)

	assert := assert.New(t)

	feed := &Feed{ID: 1, Name: "test"}

	ufeed, response, err := client.Feed.Update("test", feed)

	assert.Nil(err)
	assert.NotNil(ufeed)
	assert.NotNil(response)

	assert.Equal(1, ufeed.ID)
	assert.Equal("test-2", ufeed.Name)
	assert.NotEqual(&feed, &ufeed)
}

func TestFeedDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
		},
	)

	assert := assert.New(t)

	response, err := client.Feed.Delete("test")

	assert.Nil(err)
	assert.NotNil(response)

	assert.Equal(200, response.StatusCode)
}
