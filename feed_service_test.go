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

	assert.Equal(int64(1), feed.GetID())
	assert.Equal("feed", feed.GetName())
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

	assert.Equal(int64(1), feed.GetID())
	assert.Equal("test", feed.GetName())
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

	nfeed := &Feed{Name: Ptr("test")}

	feed, response, err := client.Feed.Create(context.Background(), nfeed)

	assert.Nil(err)
	assert.NotNil(feed)
	assert.NotNil(response)

	assert.Equal(int64(1), feed.GetID())
	assert.Equal("test", feed.GetName())
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

	feed := &Feed{ID: Ptr(int64(1)), Name: Ptr("test")}

	ufeed, response, err := client.Feed.Update(context.Background(), "test", feed)

	assert.Nil(err)
	assert.NotNil(ufeed)
	assert.NotNil(response)

	assert.Equal(int64(1), ufeed.GetID())
	assert.Equal("test-2", ufeed.GetName())
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

func TestFeedUpdateSendsZeroValues(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PATCH")
			testBody(t, r, `{"description":"","history":false}`+"\n")
			fmt.Fprint(w, `{"id":1, "key":"test", "description":"", "history":false}`)
		},
	)

	assert := assert.New(t)

	// Description and History are set to their zero values, so they must be
	// sent. Every other field is nil, so it must be left out of the request.
	ufeed, _, err := client.Feed.Update(context.Background(), "test", &Feed{
		Description: Ptr(""),
		History:     Ptr(false),
	})

	assert.Nil(err)
	assert.Equal(Ptr(""), ufeed.Description)
	assert.Equal(Ptr(false), ufeed.History)
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
	assert.Equal(Ptr(true), feed.History)
	assert.Equal(Ptr(false), feed.Enabled)
	assert.Equal(Ptr(false), feed.StatusNotify)
	assert.Equal(4320, feed.GetStatusTimeout())
	assert.True(feed.GetWritable())
	if assert.NotNil(feed.Group) {
		assert.Equal("default", feed.Group.GetKey())
		assert.Equal(int64(2), feed.Group.GetUserID())
	}
	if assert.Len(feed.Groups, 1) {
		assert.Equal("default", feed.Groups[0].GetKey())
	}
}
