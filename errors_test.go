package adafruitio

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponseBodies(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		body         string
		wantMessage  string
		wantMessages []string
	}{
		{
			"string error",
			422,
			`{"error":"request failed - Record invalid. Failed to save data to jwcooper/feeds/probe, data missing required value"}`,
			"request failed - Record invalid. Failed to save data to jwcooper/feeds/probe, data missing required value",
			[]string{"request failed - Record invalid. Failed to save data to jwcooper/feeds/probe, data missing required value"},
		},
		{
			"array of errors",
			400,
			`{"error":["Name can't be blank","Key can't be blank"],"status":"bad_request"}`,
			"Name can't be blank; Key can't be blank",
			[]string{"Name can't be blank", "Key can't be blank"},
		},
		{
			"plain text",
			500,
			`Internal Server Error`,
			"Internal Server Error",
			nil,
		},
		{
			"JSON without an error",
			500,
			`{"status":"oops"}`,
			`{"status":"oops"}`,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc(serverPattern("feeds"),
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.status)
					w.Write([]byte(tt.body))
				},
			)

			_, response, err := client.Feed.Create(context.Background(), &Feed{})

			var apiErr *ErrorResponse
			if assert.True(t, errors.As(err, &apiErr)) {
				assert.Equal(t, tt.wantMessage, apiErr.Message)
				assert.Equal(t, tt.status, apiErr.Response.StatusCode)
				if tt.wantMessages != nil {
					assert.Equal(t, tt.wantMessages, apiErr.AIOError.Messages)
				} else {
					assert.Nil(t, apiErr.AIOError)
				}
			}
			assert.Equal(t, tt.status, response.StatusCode)

			var rateErr *RateLimitError
			assert.False(t, errors.As(err, &rateErr), "only 429 is a RateLimitError")
		})
	}
}

func TestRateLimitError(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantMessage string
	}{
		// the data rate limit returns JSON
		{"data rate limit", `{"error":"request failed - Throttle error. 61 operations requested, 60 remaining"}`, "request failed - Throttle error. 61 operations requested, 60 remaining"},
		// the overall HTTP request limit returns plain text
		{"request limit", `Retry later`, "Retry later"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc(serverPattern("feeds/temperature/data/batch"),
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte(tt.body))
				},
			)

			_, response, err := client.Data.CreateBatch(context.Background(), "temperature", []*Data{{Value: Ptr("1")}})

			var rateErr *RateLimitError
			if assert.True(t, errors.As(err, &rateErr), "got %T", err) {
				assert.Equal(t, tt.wantMessage, rateErr.Message)
				assert.Contains(t, rateErr.Error(), "429")
			}

			// a RateLimitError is also an ErrorResponse
			var apiErr *ErrorResponse
			if assert.True(t, errors.As(err, &apiErr)) {
				assert.Equal(t, http.StatusTooManyRequests, apiErr.Response.StatusCode)
			}
			assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
		})
	}
}
