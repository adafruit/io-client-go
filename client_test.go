package adafruitio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testUser = "test_username"

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Adafruit IO client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient(testUser, "test-key")
	client.SetBaseURL(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func serverPattern(pattern string) string {
	return fmt.Sprintf("%s/%s/%s", APIPath, testUser, pattern)
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %s, want %s", header, got, want)
	}
}

func testQuery(t *testing.T, r *http.Request, field string, want string) {
	r.ParseForm()

	if len(r.Form[field]) < 1 {
		t.Errorf("expected Form[%q] to have a value", field)
		fmt.Printf("FORM: %q\n", r.Form)
		return
	}

	if got := r.Form[field][0]; got != want {
		t.Errorf("Form[%q] returned %s, want %s", field, got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

func TestClientInitiation(t *testing.T) {
	assert := assert.New(t)

	c := NewClient("GIVEN USER", "GIVEN KEY")
	u, k := c.GetUserKey()

	assert.Equal("GIVEN USER", u, "expected to find GIVEN USER")
	assert.Equal("GIVEN KEY", k, "expected to find GIVEN KEY")
}

func TestClientAuthentication(t *testing.T) {
	setup()
	defer teardown()
	assert := assert.New(t)

	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			testHeader(t, r, "X-AIO-Key", "test-key")
			fmt.Fprintf(w, "ok")
		},
	)

	req, err := client.NewRequest("GET", "/", nil)
	assert.Nil(err)
	assert.NotNil(req)

	resp, err := client.Do(context.Background(), req, nil)
	assert.Nil(err)
	assert.NotNil(resp)
}

func TestCheckResponseNullBody(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://io.adafruit.com/api/v2/test", nil)
	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("null")),
		Request:    req,
	}

	err := CheckResponse(resp)

	if assert.IsType(&ErrorResponse{}, err) {
		assert.Equal("null", err.(*ErrorResponse).Message)
	}
}

// countingTransport counts requests before passing them to the default transport.
type countingTransport struct {
	count int
}

func (t *countingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.count++
	return http.DefaultTransport.RoundTrip(r)
}

func TestWithHTTPClient(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds"),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[]`)
		},
	)

	transport := &countingTransport{}
	c := NewClient(testUser, "test-key", WithHTTPClient(&http.Client{Transport: transport}))
	c.SetBaseURL(server.URL)

	_, _, err := c.Feed.All(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, 1, transport.count)
}

func TestWithHTTPClientNil(t *testing.T) {
	c := NewClient(testUser, "test-key", WithHTTPClient(nil))

	assert.Equal(t, http.DefaultClient, c.client)
}

func TestDoCanceledContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds"),
		func(w http.ResponseWriter, r *http.Request) {
			t.Error("request should not reach the server")
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	feeds, response, err := client.Feed.All(ctx)

	assert.True(t, errors.Is(err, context.Canceled), "got error %v", err)
	assert.Nil(t, feeds)
	assert.Nil(t, response)
}

func TestDoContextDeadline(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds"),
		func(w http.ResponseWriter, r *http.Request) {
			// respond only after the client has given up
			select {
			case <-r.Context().Done():
			case <-time.After(5 * time.Second):
			}
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, _, err := client.Feed.All(ctx)

	assert.True(t, errors.Is(err, context.DeadlineExceeded), "got error %v", err)
}

func TestDoNilContext(t *testing.T) {
	setup()
	defer teardown()

	req, err := client.NewRequest("GET", "feeds", nil)
	assert.Nil(t, err)

	//lint:ignore SA1012 testing that a nil context is rejected
	response, err := client.Do(nil, req, nil)

	assert.Equal(t, errNilContext, err)
	assert.Nil(t, response)
}

func TestNextLink(t *testing.T) {
	const self = "https://io.adafruit.com/api/v2/u/feeds/f/data?end_time=2026-11-24+17%3A51%3A16+UTC&limit=2"
	const next = "https://io.adafruit.com/api/v2/u/feeds/f/data?before=abc&end_time=2026-11-24+17%3A51%3A16+UTC&limit=2"

	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"Adafruit IO middle page", "<" + self + `>;rel="next", <` + next + ">", next},
		{"Adafruit IO last page", "<" + self + ">;", ""},
		{"URL with a comma", "<" + self + `&include=value,created_at>;rel="next", <` + next + "&include=value,created_at>", next + "&include=value,created_at"},
		{"empty", "", ""},
		{"standard, next only", "<" + next + `>; rel="next"`, next},
		{"standard, with prev and last", "<" + self + `>; rel="prev", <` + next + `>; rel="next", <` + self + `>; rel="last"`, next},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, nextLink(tt.header), tt.name)
	}
}

func TestResponsePagination(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/temperature/data"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Pagination-Total", "5")
			w.Header().Set("Link", `<http://example.com/data?limit=2>;rel="next", <http://example.com/data?before=abc%3D%3D&limit=2>`)
			fmt.Fprint(w, `[{"id":"5"},{"id":"4"}]`)
		},
	)

	_, response, err := client.Data.All(context.Background(), "temperature", nil)

	assert.Nil(t, err)
	assert.Equal(t, 5, response.Total)
	assert.Equal(t, "abc==", response.NextBefore)
}

func TestResponseNotPaginated(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(serverPattern("feeds/temperature"),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"key":"temperature"}`)
		},
	)

	_, response, err := client.Feed.Get(context.Background(), "temperature")

	assert.Nil(t, err)
	assert.Equal(t, 0, response.Total)
	assert.Equal(t, "", response.NextBefore)
}
