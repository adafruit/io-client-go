// Portions of code in this file are adapted from the go-github
// project located at https://github.com/google/go-github

package adafruitio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	BaseURL       = "https://io.adafruit.com"
	xAIOKeyHeader = "X-AIO-Key"
)

type Client struct {
	// Base HTTP client used to talk to io.adafruit.com
	client *http.Client

	// Base URL for API requests. Defaults to public adafruit io URL.
	BaseURL *url.URL

	APIKey    string
	UserAgent string

	// Services that make up adafruit io.
	Data  *DataService
	Feed  *FeedService
	Group *GroupService
}

// Response wraps http.Response and adds fields unique to Adafruit's API.
type Response struct {
	*http.Response
}

func (r *Response) Debug() {
	all, _ := ioutil.ReadAll(r.Body)
	fmt.Println("---")
	fmt.Println(string(all))
	fmt.Println("---")
}

type AIOError struct {
	Message string `json:"error"`
}

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that carried the error message
	Message  string
	AIOError *AIOError
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"%v %v %v: %v",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
		r.Message,
	)
}

func NewClient(key string) *Client {
	c := &Client{APIKey: key}

	c.BaseURL, _ = url.Parse(BaseURL)
	c.UserAgent = "Adafruit IO Go Client v0.1"

	c.client = http.DefaultClient

	c.Data = &DataService{client: c}
	c.Feed = &FeedService{client: c}
	c.Group = &GroupService{client: c}

	return c
}

// SetFeed takes a Feed record as a parameter and uses that feed for all
// subsequent Data related API calls.
//
// A Feed must be set before making calls to the Data service.
func (c *Client) SetFeed(feed *Feed) {
	c.Feed.CurrentFeed = feed
}

func (c *Client) checkFeed() error {
	if c.Feed.CurrentFeed == nil {
		return fmt.Errorf("CurrentFeed must be set")
	}
	return nil
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
//
// adapted from https://github.com/google/go-github
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}

	// read response body into Error.Message
	body, _ := ioutil.ReadAll(r.Body)

	// try to unmarshal error response Body into AIOError record
	jerr := json.Unmarshal(body, &errorResponse.AIOError)
	if jerr != nil {
		fmt.Println("> failed to parse response body as JSON")
		// failed to unmarhsal API Error, use body as Message
		errorResponse.Message = string(body)
	} else {
		fmt.Println("> parsed response body as JSON", errorResponse.AIOError)
		errorResponse.Message = errorResponse.AIOError.Message
	}

	return errorResponse
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
//
// adapted from https://github.com/google/go-github
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Authentication v1
	req.Header.Add(xAIOKeyHeader, c.APIKey)

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// adapted from https://github.com/google/go-github
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := &Response{
		Response: resp,
	}

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}
