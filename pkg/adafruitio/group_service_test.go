package adafruitio

import (
	"fmt"
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

	groups, response, err := client.Group.All()

	assert.Nil(err)
	assert.NotNil(groups)
	assert.NotNil(response)

	group := groups[0]

	assert.Equal(1, group.ID)
	assert.Equal("group", group.Name)
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

	group, response, err := client.Group.Get("test")

	assert.Nil(err)
	assert.NotNil(group)
	assert.NotNil(response)

	assert.Equal(1, group.ID)
	assert.Equal("test", group.Name)
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

	ngroup := &Group{Name: "test"}

	group, response, err := client.Group.Create(ngroup)

	assert.Nil(err)
	assert.NotNil(group)
	assert.NotNil(response)

	assert.Equal(1, group.ID)
	assert.Equal("test", group.Name)
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

	group := &Group{ID: 1, Name: "test"}

	ugroup, response, err := client.Group.Update("test", group)

	assert.Nil(err)
	assert.NotNil(ugroup)
	assert.NotNil(response)

	assert.Equal(1, ugroup.ID)
	assert.Equal("test-2", ugroup.Name)
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

	response, err := client.Group.Delete("test")

	assert.Nil(err)
	assert.NotNil(response)

	assert.Equal(200, response.StatusCode)
}
