package adafruitio

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `[
				{
					"id": 1,
					"name": "group",
					"key": "group",
					"description": null,
					"source": null,
					"properties": null,
					"source_keys": [
						"temperature",
						"pressure",
						"humidity",
						"temp_min",
						"temp_max",
						"wind_speed",
						"wind_direction",
						"rainfall",
						"name"
					],
					"created_at": "2016-05-26T18:50:09.695Z",
					"updated_at": "2016-05-27T15:08:11.661Z",
					"feeds": [
						{
							"id": 1,
							"name": "beta-test",
							"key": "beta-test",
							"unit_type": null,
							"unit_symbol": null,
							"mode": null,
							"history": true,
							"last_value": "67.123441",
							"last_value_at": "2016-05-26T14:37:14.306Z",
							"created_at": "2016-05-23T18:01:39.753Z",
							"updated_at": "2016-05-27T15:08:11.663Z",
							"stream": {
								"id": 565778728,
								"value": "67.123441",
								"lat": null,
								"lon": null,
								"ele": null,
								"completed_at": null,
								"created_at": "2016-05-26T14:37:14.306Z"
							}
						}
					]
				}
			]`)
		},
	)

	assert := assert.New(t)

	groups, response, err := client.Group.All(context.Background())

	assert.Nil(err)
	assert.NotNil(groups)
	assert.NotNil(response)

	group := groups[0]

	assert.Equal(int64(1), group.GetID())
	assert.Equal("group", group.GetName())
}

func TestGroupGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":1, "name":"test", "feeds": [{"id": 1}]}`)
		},
	)

	assert := assert.New(t)

	group, response, err := client.Group.Get(context.Background(), "test")

	assert.Nil(err)
	assert.NotNil(group)
	assert.NotNil(response)

	assert.Equal(int64(1), group.GetID())
	assert.Equal("test", group.GetName())
}

func TestGroupCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":1, "name":"test"}`)
		},
	)

	assert := assert.New(t)

	ngroup := &Group{Name: Ptr("test")}

	group, response, err := client.Group.Create(context.Background(), ngroup)

	assert.Nil(err)
	assert.NotNil(group)
	assert.NotNil(response)

	assert.Equal(int64(1), group.GetID())
	assert.Equal("test", group.GetName())
}

func TestGroupUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PATCH")
			fmt.Fprint(w, `{"id":1, "name":"test-2"}`)
		},
	)

	assert := assert.New(t)

	group := &Group{ID: Ptr(int64(1)), Name: Ptr("test")}

	ugroup, response, err := client.Group.Update(context.Background(), "test", group)

	assert.Nil(err)
	assert.NotNil(ugroup)
	assert.NotNil(response)

	assert.Equal(int64(1), ugroup.GetID())
	assert.Equal("test-2", ugroup.GetName())
	assert.NotEqual(&group, &ugroup)
}

func TestGroupDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
		},
	)

	assert := assert.New(t)

	response, err := client.Group.Delete(context.Background(), "test")

	assert.Nil(err)
	assert.NotNil(response)

	assert.Equal(200, response.StatusCode)
}

func TestGroupCreateOmitsEmptyVisibility(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			body, _ := io.ReadAll(r.Body)
			assert.NotContains(t, string(body), "visibility")
			fmt.Fprint(w, `{"id":1, "name":"test"}`)
		},
	)

	_, _, err := client.Group.Create(context.Background(), &Group{Name: Ptr("test")})

	assert.Nil(t, err)
}

func TestGroupDeleteError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/test"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		},
	)

	assert := assert.New(t)

	response, err := client.Group.Delete(context.Background(), "test")

	assert.NotNil(err)
	if assert.NotNil(response) {
		assert.Equal(404, response.StatusCode)
	}
}

func TestGroupFeeds(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/weather/feeds"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `[{"id": 1, "key": "weather.temperature"}, {"id": 2, "key": "weather.humidity"}]`)
		},
	)

	feeds, _, err := client.Group.Feeds(context.Background(), "weather")

	assert.Nil(t, err)
	if assert.Len(t, feeds, 2) {
		assert.Equal(t, "weather.temperature", feeds[0].GetKey())
	}
}

func TestGroupCreateFeed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/weather/feeds"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testBody(t, r, `{"name":"Temperature"}`+"\n")
			fmt.Fprint(w, `{"id": 1, "name": "Temperature", "key": "weather.temperature"}`)
		},
	)

	feed, _, err := client.Group.CreateFeed(context.Background(), "weather", &Feed{Name: Ptr("Temperature")})

	assert.Nil(t, err)
	assert.Equal(t, "weather.temperature", feed.GetKey())
}

func TestGroupAddAndRemoveFeed(t *testing.T) {
	setup()
	defer teardown()

	for _, action := range []string{"add", "remove"} {
		action := action
		mux.HandleFunc(serverPattern("groups/weather/"+action),
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				testBody(t, r, `{"feed_key":"temperature"}`+"\n")
				fmt.Fprintf(w, `{"id": 1, "key": "weather", "description": "%s"}`, action)
			},
		)
	}

	group, _, err := client.Group.AddFeed(context.Background(), "weather", "temperature")
	assert.Nil(t, err)
	assert.Equal(t, "add", group.GetDescription())

	group, _, err = client.Group.RemoveFeed(context.Background(), "weather", "temperature")
	assert.Nil(t, err)
	assert.Equal(t, "remove", group.GetDescription())
}

func TestGroupLastValues(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("groups/weather/data/last_value"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testQuery(t, r, "type", "json")
			fmt.Fprint(w, `{"feeds": {"weather.temperature": "22.5", "weather.humidity": null}}`)
		},
	)

	values, _, err := client.Group.LastValues(context.Background(), "weather")

	assert.Nil(t, err)
	assert.Equal(t, map[string]*string{
		"weather.temperature": Ptr("22.5"),
		"weather.humidity":    nil,
	}, values)
}
