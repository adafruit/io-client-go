package adafruitio

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestData_EmptyFeedKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("request should not reach the server: %s %s", r.Method, r.URL)
		},
	)

	ctx := context.Background()
	calls := map[string]func() error{
		"All":         func() error { _, _, err := client.Data.All(ctx, "", nil); return err },
		"Search":      func() error { _, _, err := client.Data.Search(ctx, "", nil); return err },
		"Get":         func() error { _, _, err := client.Data.Get(ctx, "", "1"); return err },
		"Create":      func() error { _, _, err := client.Data.Create(ctx, "", &Data{}); return err },
		"CreateBatch": func() error { _, _, err := client.Data.CreateBatch(ctx, "", nil); return err },
		"Update":      func() error { _, _, err := client.Data.Update(ctx, "", "1", &Data{}); return err },
		"Delete":      func() error { _, err := client.Data.Delete(ctx, "", "1"); return err },
		"Next":        func() error { _, _, err := client.Data.Next(ctx, ""); return err },
		"Prev":        func() error { _, _, err := client.Data.Prev(ctx, ""); return err },
		"First":       func() error { _, _, err := client.Data.First(ctx, ""); return err },
		"Last":        func() error { _, _, err := client.Data.Last(ctx, ""); return err },
	}

	for name, call := range calls {
		assert.Equal(t, errEmptyFeedKey, call(), name)
	}
}

func TestDataCreate(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{"id":"1", "value":"67.112"}`)
		},
	)

	assert := assert.New(t)

	val := "67.112"

	dp := &Data{
		Value: Ptr(val),
	}
	datapoint, response, err := client.Data.Create(context.Background(), "temperature", dp)

	assert.Nil(err)
	assert.NotNil(datapoint)
	assert.NotNil(response)

	assert.Equal("1", datapoint.GetID())
	assert.Equal(val, datapoint.GetValue())
}

func TestDataGet(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc(serverPattern("feeds/temperature/data/1"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":"1", "value":"67.112"}`)
		},
	)

	assert := assert.New(t)

	datapoint, response, err := client.Data.Get(context.Background(), "temperature", "1")

	assert.Nil(err)
	assert.NotNil(datapoint)
	assert.NotNil(response)

	assert.Equal("1", datapoint.GetID())
	assert.Equal("67.112", datapoint.GetValue())
}

func TestAllDataNoFilter(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `[{"id":"1", "value":"67.112"}]`)
		},
	)

	assert := assert.New(t)

	// with no params
	datapoints, response, err := client.Data.All(context.Background(), "temperature", nil)
	datapoint := datapoints[0]

	assert.Nil(err)
	assert.NotNil(datapoint)
	assert.NotNil(response)

	assert.Equal("1", datapoint.GetID())
	assert.Equal("67.112", datapoint.GetValue())
}

func TestAllDataFilter(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testQuery(t, r, "start_time", "2000-01-01T00:00:00Z")
			testQuery(t, r, "end_time", "2010-01-01T00:00:00Z")
			fmt.Fprint(w, `[{"id":"1", "value":"67.112"}]`)
		},
	)

	assert := assert.New(t)

	// with no params
	datapoints, response, err := client.Data.All(context.Background(), "temperature", &DataFilter{
		StartTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	datapoint := datapoints[0]

	assert.Nil(err)
	assert.NotNil(datapoint)
	assert.NotNil(response)

	assert.Equal("1", datapoint.GetID())
	assert.Equal("67.112", datapoint.GetValue())
}

func TestSearchDataFilter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testQuery(t, r, "start_time", "2000-01-01T00:00:00Z")
			testQuery(t, r, "end_time", "2010-01-01T00:00:00Z")
			fmt.Fprint(w, `[{"id":"1", "value":"67.112"}]`)
		},
	)

	assert := assert.New(t)

	datapoints, response, err := client.Data.Search(context.Background(), "temperature", &DataFilter{
		StartTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	assert.Nil(err)
	assert.NotNil(response)
	if assert.Len(datapoints, 1) {
		assert.Equal("1", datapoints[0].GetID())
		assert.Equal("67.112", datapoints[0].GetValue())
	}
}

func TestDataDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test/data/1"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
		},
	)

	assert := assert.New(t)

	response, err := client.Data.Delete(context.Background(), "test", "1")

	assert.Nil(err)
	assert.NotNil(response)

	assert.Equal(200, response.StatusCode)
}

func TestDataQueue(t *testing.T) {
	setup()
	defer teardown()

	// prepare endpoint URL for just this request
	mux.HandleFunc(serverPattern("feeds/temperature/data/next"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":"1", "value":"1"}`)
		},
	)

	mux.HandleFunc(serverPattern("feeds/temperature/data/previous"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":"2", "value":"2"}`)
		},
	)

	mux.HandleFunc(serverPattern("feeds/temperature/data/last"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":"3", "value":"3"}`)
		},
	)
	mux.HandleFunc(serverPattern("feeds/temperature/data/first"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"id":"1", "value":"1"}`)
		},
	)
	assert := assert.New(t)

	var (
		datapoint *Data
		response  *Response
		err       error
	)

	datapoint, response, err = client.Data.Next(context.Background(), "temperature")
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal("1", datapoint.GetID())
	assert.Equal("1", datapoint.GetValue())

	datapoint, response, err = client.Data.Prev(context.Background(), "temperature")
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal("2", datapoint.GetID())
	assert.Equal("2", datapoint.GetValue())

	datapoint, response, err = client.Data.Last(context.Background(), "temperature")
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal("3", datapoint.GetID())
	assert.Equal("3", datapoint.GetValue())

	datapoint, response, err = client.Data.First(context.Background(), "temperature")
	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal("1", datapoint.GetID())
	assert.Equal("1", datapoint.GetValue())

}

func TestDataDeleteError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/test/data/1"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		},
	)

	assert := assert.New(t)

	response, err := client.Data.Delete(context.Background(), "test", "1")

	assert.NotNil(err)
	if assert.NotNil(response) {
		assert.Equal(404, response.StatusCode)
	}
}

func TestAllDataPagingFilter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testQuery(t, r, "limit", "2")
			testQuery(t, r, "include", "value,created_at")
			testQuery(t, r, "before", "abc==")
			fmt.Fprint(w, `[]`)
		},
	)

	_, _, err := client.Data.All(context.Background(), "temperature", &DataFilter{
		Limit:   2,
		Include: []string{"value", "created_at"},
		Before:  "abc==",
	})

	assert.Nil(t, err)
}

func TestAllDataPaging(t *testing.T) {
	setup()
	defer teardown()

	// three pages, linked by before cursors, in the API's Link header format
	pages := map[string]struct {
		body string
		next string
	}{
		"":   {`[{"value":"5"},{"value":"4"}]`, "p2"},
		"p2": {`[{"value":"3"},{"value":"2"}]`, "p3"},
		"p3": {`[{"value":"1"}]`, ""},
	}

	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			page, ok := pages[r.URL.Query().Get("before")]
			if !ok {
				t.Errorf("unexpected before cursor %q", r.URL.Query().Get("before"))
				return
			}
			self := "<http://example.com/data?limit=2>;"
			if page.next != "" {
				w.Header().Set("Link", self+`rel="next", <http://example.com/data?before=`+page.next+`&limit=2>`)
			} else {
				w.Header().Set("Link", self)
			}
			w.Header().Set("X-Pagination-Total", "5")
			fmt.Fprint(w, page.body)
		},
	)

	var values []string
	filter := &DataFilter{Limit: 2}
	for {
		page, response, err := client.Data.All(context.Background(), "temperature", filter)
		if !assert.Nil(t, err) {
			return
		}
		for _, d := range page {
			values = append(values, d.GetValue())
		}
		if response.NextBefore == "" {
			break
		}
		filter.Before = response.NextBefore
	}

	assert.Equal(t, []string{"5", "4", "3", "2", "1"}, values)
}

func TestDataCreateBatch(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/temperature/data/batch"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testBody(t, r, `{"data":[{"value":"1"},{"value":"2"}]}`+"\n")
			fmt.Fprint(w, `[{"id":"a", "value":"1"}, {"id":"b", "value":"2"}]`)
		},
	)

	assert := assert.New(t)

	created, response, err := client.Data.CreateBatch(context.Background(), "temperature", []*Data{
		{Value: Ptr("1")},
		{Value: Ptr("2")},
	})

	assert.Nil(err)
	assert.NotNil(response)
	if assert.Len(created, 2) {
		assert.Equal("a", created[0].GetID())
		assert.Equal("2", created[1].GetValue())
	}
}
