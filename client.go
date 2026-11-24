// Portions of code in this file are adapted from the go-github
// project located at https://github.com/google/go-github

package adafruitio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	BaseURL       = "https://io.adafruit.com"
	APIPath       = "/api/v2"
	xAIOKeyHeader = "X-AIO-Key"
)

var errNilContext = errors.New("context must be non-nil")

type Client struct {
	// Base HTTP client used to talk to io.adafruit.com
	client *http.Client

	// Base URL for API requests. Defaults to public adafruit io URL.
	baseURL *url.URL

	apiKey    string
	username  string
	userAgent string

	// Services that make up adafruit io.
	Data  *DataService
	Feed  *FeedService
	Group *GroupService
	User  *UserService
}

// Response wraps http.Response and adds fields unique to Adafruit's API.
type Response struct {
	*http.Response

	// Total is the total number of records matching a paginated request,
	// across all pages, from the X-Pagination-Total header. The API notes it
	// may be up to 5 minutes behind. It is 0 for requests that aren't
	// paginated.
	Total int

	// NextBefore is the cursor for the next page of a paginated request, such
	// as DataService.All. Set it as the Before field of the next request's
	// DataFilter. It is empty when there are no more pages.
	NextBefore string
}

// newResponse wraps resp and populates the pagination fields from its headers.
func newResponse(resp *http.Response) *Response {
	response := &Response{Response: resp}
	response.Total, _ = strconv.Atoi(resp.Header.Get("X-Pagination-Total"))
	if next := nextLink(resp.Header.Get("Link")); next != "" {
		if u, err := url.Parse(next); err == nil {
			response.NextBefore = u.Query().Get("before")
		}
	}
	return response
}

// nextLink returns the URL of the next page from a Link header, or "" if
// there isn't one.
//
// Adafruit IO's Link header doesn't follow the usual format. It looks like
//
//	<current page URL>;rel="next", <next page URL>
//
// where rel="next" marks the URL that follows it, rather than the one before
// it. On the last page, only the current page URL is present. To stay
// compatible if the API adopts the standard format, where the URL before
// rel="next" is the next page, that URL is used when nothing follows it.
func nextLink(header string) string {
	links := parseLinks(header)
	for i, link := range links {
		if !strings.Contains(link.params, `rel="next"`) {
			continue
		}
		if i+1 < len(links) && !strings.Contains(links[i+1].params, "rel=") {
			return links[i+1].url
		}
		return link.url
	}
	return ""
}

// link is one entry of a Link header: a URL and the parameters after it.
type link struct {
	url    string
	params string
}

// parseLinks splits a Link header into its entries. It finds each URL by its
// angle brackets rather than splitting on commas, because URLs may contain
// commas, for example in an include=value,created_at query parameter.
func parseLinks(header string) []link {
	var links []link
	for {
		start := strings.Index(header, "<")
		if start < 0 {
			return links
		}
		end := strings.Index(header[start:], ">")
		if end < 0 {
			return links
		}
		end += start

		params := header[end+1:]
		if next := strings.Index(params, "<"); next >= 0 {
			params = params[:next]
		}
		links = append(links, link{url: header[start+1 : end], params: params})
		header = header[end+1:]
	}
}

// AIOError is the body of an API error response.
type AIOError struct {
	// Messages are the error messages. Most errors have one, but validation
	// errors, such as creating a feed without a name, may have several.
	Messages []string

	// Message is Messages joined with "; ".
	Message string
}

// UnmarshalJSON implements the json.Unmarshaler interface. The API returns
// the error as either a string or an array of strings:
//
//	{"error": "request failed - cannot process request"}
//	{"error": ["Name can't be blank", "Key can't be blank"], "status": "bad_request"}
func (e *AIOError) UnmarshalJSON(b []byte) error {
	var body struct {
		Error json.RawMessage `json:"error"`
	}
	if err := json.Unmarshal(b, &body); err != nil {
		return err
	}

	var one string
	var many []string
	switch {
	case len(body.Error) == 0 || string(body.Error) == "null":
		e.Messages = nil
	case json.Unmarshal(body.Error, &one) == nil:
		e.Messages = []string{one}
	case json.Unmarshal(body.Error, &many) == nil:
		e.Messages = many
	default:
		return fmt.Errorf("unexpected error field: %s", body.Error)
	}

	e.Message = strings.Join(e.Messages, "; ")
	return nil
}

// ErrorResponse reports an error returned by the API: any response with a
// status code outside the 200 range.
type ErrorResponse struct {
	Response *http.Response // HTTP response that carried the error message

	// Message is the API's error message, or the response body if it
	// doesn't contain one.
	Message string

	// AIOError is the decoded error body, or nil if the body isn't a JSON
	// error.
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

// RateLimitError is returned when the API responds with 429 Too Many
// Requests: either the account's data rate limit (data points created,
// updated, or deleted per minute) or the overall limit on HTTP requests per
// minute was exceeded. The API doesn't say when to retry; its documentation
// suggests waiting a minute.
//
// RateLimitError wraps an *ErrorResponse, so errors.As finds either type.
type RateLimitError struct {
	*ErrorResponse
}

// Unwrap returns the underlying *ErrorResponse.
func (r *RateLimitError) Unwrap() error {
	return r.ErrorResponse
}

// ClientOption configures a Client. Pass options to NewClient.
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client used to make API requests, for example
// to configure timeouts, proxies, or transport settings. A nil client is
// ignored. By default, http.DefaultClient is used.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.client = httpClient
		}
	}
}

// NewClient returns a Client for the Adafruit IO account identified by
// username, authenticating with the given API key.
func NewClient(username, key string, opts ...ClientOption) *Client {
	c := &Client{username: username, apiKey: key}

	c.SetBaseURL(BaseURL)
	c.userAgent = fmt.Sprintf("AdafruitIO-Go/%v (%v %v)", Version, runtime.GOOS, runtime.Version())

	c.client = http.DefaultClient

	c.Data = &DataService{client: c}
	c.Feed = &FeedService{client: c}
	c.Group = &GroupService{client: c}
	c.User = &UserService{client: c}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetBaseURL updates the base URL to use. Mainly here for use in unit testing
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL, _ = url.Parse(fmt.Sprintf("%s%s/%s/", baseURL, APIPath, c.username))
}

func (c *Client) GetUserKey() (username string, apikey string) {
	return c.username, c.apiKey
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
	body, _ := io.ReadAll(r.Body)

	// try to unmarshal error response Body into AIOError record
	jerr := json.Unmarshal(body, &errorResponse.AIOError)
	if jerr != nil || errorResponse.AIOError == nil {
		// body was not a JSON error object, use body as Message
		errorResponse.AIOError = nil
		errorResponse.Message = string(body)
	} else if errorResponse.AIOError.Message == "" {
		// JSON without an error message, use body as Message
		errorResponse.AIOError = nil
		errorResponse.Message = string(body)
	} else {
		errorResponse.Message = errorResponse.AIOError.Message
	}

	if r.StatusCode == http.StatusTooManyRequests {
		return &RateLimitError{ErrorResponse: errorResponse}
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

	u := c.baseURL.ResolveReference(rel)

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
	req.Header.Add(xAIOKeyHeader, c.apiKey)

	if c.userAgent != "" {
		req.Header.Add("User-Agent", c.userAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The request is bound to ctx. If ctx is cancelled or its deadline passes,
// Do returns ctx.Err().
//
// adapted from https://github.com/google/go-github
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errNilContext
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		// If the context was cancelled, its error is more useful than the
		// wrapped transport error.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(io.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := newResponse(resp)

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

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
