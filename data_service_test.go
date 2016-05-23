package adafruitio_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	aio "github.com/adafruit/io-client-go"
)

func TestData_MissingFeed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/feeds/temperature/data",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":1, "value":"67.112"}`)
		},
	)

	assert := assert.New(t)

	dp := &aio.DataPoint{}
	datapoint, response, err := client.Data.Create(dp)

	assert.NotNil(err)
	assert.Nil(datapoint)
	assert.Nil(response)

	assert.Equal(err.Error(), "Feed.Name cannot be empty")
}

func TestData_Unauthenticated(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/feeds/temperature/data",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":1, "value":"67.112"}`)
		},
	)

	assert := assert.New(t)

	dp := &aio.DataPoint{}
	datapoint, response, err := client.Data.Create(dp)

	assert.NotNil(err)
	assert.Nil(datapoint)
	assert.Nil(response)

	assert.Equal(err.Error(), "Feed.Name cannot be empty")
}

func TestDataCreate(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc("/feeds/temperature/data",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":1, "value":"67.112"}`)
		},
	)

	assert := assert.New(t)

	client.SetFeed("temperature")

	val := json.Number("67.112")

	dp := &aio.DataPoint{
		Value: &val,
	}
	datapoint, response, err := client.Data.Create(dp)

	assert.Nil(err)
	assert.NotNil(datapoint)
	assert.NotNil(response)

	assert.Equal(int64(1), *datapoint.ID)
	assert.Equal(val, *datapoint.Value)
}

func TestDataSend(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc("/feeds/temperature/data/send",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":1, "value":"67.112"}`)
		},
	)

	assert := assert.New(t)

	client.SetFeed("temperature")

	val := json.Number("67.112")

	dp := &aio.DataPoint{
		Value: &val,
	}
	datapoint, response, err := client.Data.Send(dp)

	assert.Nil(err)
	assert.NotNil(datapoint)
	assert.NotNil(response)

	assert.Equal(int64(1), *datapoint.ID)
	assert.Equal(val, *datapoint.Value)
}
