package adafruitio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
	b, err := ioutil.ReadAll(r.Body)
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

	resp, err := client.Do(req, nil)
	assert.Nil(err)
	assert.NotNil(resp)
}
