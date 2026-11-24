package adafruitio

import (
	"context"
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

	feeds, response, err := client.Feed.All(context.Background())

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

	feed, response, err := client.Feed.Get(context.Background(), "test")

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

	feed, response, err := client.Feed.Create(context.Background(), nfeed)

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

	ufeed, response, err := client.Feed.Update(context.Background(), "test", feed)

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

	response, err := client.Feed.Delete(context.Background(), "test")

	assert.Nil(err)
	assert.NotNil(response)

	assert.Equal(200, response.StatusCode)
}

func TestFeedDeleteError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		},
	)

	assert := assert.New(t)

	response, err := client.Feed.Delete(context.Background(), "test")

	assert.NotNil(err)
	if assert.NotNil(response) {
		assert.Equal(404, response.StatusCode)
	}
}

func TestFeedUpdateSendsFalse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PATCH")
			testBody(t, r, `{"history":false}`+"\n")
			fmt.Fprint(w, `{"id":1, "key":"test", "history":false}`)
		},
	)

	assert := assert.New(t)

	// History is explicitly false; Enabled and StatusNotify are nil, so they
	// must be left out of the request.
	ufeed, _, err := client.Feed.Update(context.Background(), "test", &Feed{History: Bool(false)})

	assert.Nil(err)
	if assert.NotNil(ufeed.History) {
		assert.False(*ufeed.History)
	}
}

func TestFeedGetDecodesFields(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{
				"id": 1, "key": "test", "history": true, "enabled": false,
				"status_notify": false, "status_timeout": 4320, "writable": true,
				"group": {"id": 2, "key": "default", "name": "Default", "user_id": 2},
				"groups": [{"id": 2, "key": "default", "name": "Default", "user_id": 2}]
			}`)
		},
	)

	assert := assert.New(t)

	feed, _, err := client.Feed.Get(context.Background(), "test")

	assert.Nil(err)
	assert.Equal(Bool(true), feed.History)
	assert.Equal(Bool(false), feed.Enabled)
	assert.Equal(Bool(false), feed.StatusNotify)
	assert.Equal(4320, feed.StatusTimeout)
	assert.True(feed.Writable)
	if assert.NotNil(feed.Group) {
		assert.Equal("default", feed.Group.Key)
		assert.Equal(2, feed.Group.UserID)
	}
	if assert.Len(feed.Groups, 1) {
		assert.Equal("default", feed.Groups[0].Key)
	}
}
