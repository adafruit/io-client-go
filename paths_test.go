package adafruitio

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPathEscaping checks that keys and IDs are sent as single, escaped path
// segments. It uses its own server, rather than setup's ServeMux, because
// ServeMux cleans paths before handlers see them.
func TestPathEscaping(t *testing.T) {
	var requestURI string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		if strings.HasSuffix(r.URL.Path, "/data") {
			fmt.Fprint(w, `[]`) // a list endpoint
			return
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()

	c := NewClient(testUser, "test-key")
	c.SetBaseURL(server.URL)
	ctx := context.Background()
	base := APIPath + "/" + testUser + "/"

	tests := []struct {
		name string
		call func() error
		want string
	}{
		{
			"plain key",
			func() error { _, _, err := c.Feed.Get(ctx, "temperature"); return err },
			base + "feeds/temperature",
		},
		{
			"group-scoped key keeps its dot",
			func() error { _, _, err := c.Feed.Get(ctx, "weather.temperature"); return err },
			base + "feeds/weather.temperature",
		},
		{
			"feed name with a space",
			func() error { _, _, err := c.Feed.Get(ctx, "Welcome Feed"); return err },
			base + "feeds/Welcome%20Feed",
		},
		{
			"slash, question mark, and hash",
			func() error { _, _, err := c.Feed.Get(ctx, "a/b?c#d"); return err },
			base + "feeds/a%2Fb%3Fc%23d",
		},
		{
			"dot-dot key",
			func() error { _, err := c.Feed.Delete(ctx, ".."); return err },
			base + "feeds/%2E%2E",
		},
		{
			"dot key",
			func() error { _, _, err := c.Group.Get(ctx, "."); return err },
			base + "groups/%2E",
		},
		{
			"data record ID",
			func() error { _, _, err := c.Data.Get(ctx, "temperature", "../../user"); return err },
			base + "feeds/temperature/data/..%2F..%2Fuser",
		},
		{
			"query string kept after an escaped key",
			func() error { _, _, err := c.Data.LastValue(ctx, "my feed"); return err },
			base + "feeds/my%20feed/data/last_value?type=json",
		},
		{
			"group query string",
			func() error { _, _, err := c.Group.LastValues(ctx, "my group"); return err },
			base + "groups/my%20group/data/last_value?type=json",
		},
		{
			"filter query string",
			func() error { _, _, err := c.Data.All(ctx, "a?b", &DataFilter{Limit: 5}); return err },
			base + "feeds/a%3Fb/data?limit=5",
		},
	}

	for _, tt := range tests {
		requestURI = ""
		assert.NoError(t, tt.call(), tt.name)
		assert.Equal(t, tt.want, requestURI, tt.name)
	}
}

func TestEmptyKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("request should not reach the server: %s %s", r.Method, r.URL)
		},
	)

	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
		want error
	}{
		{"Feed.Get", func() error { _, _, err := client.Feed.Get(ctx, ""); return err }, errEmptyFeedKey},
		{"Feed.Update", func() error { _, _, err := client.Feed.Update(ctx, "", &Feed{}); return err }, errEmptyFeedKey},
		{"Feed.Delete", func() error { _, err := client.Feed.Delete(ctx, ""); return err }, errEmptyFeedKey},
		{"Group.Get", func() error { _, _, err := client.Group.Get(ctx, ""); return err }, errEmptyGroupKey},
		{"Group.Update", func() error { _, _, err := client.Group.Update(ctx, "", &Group{}); return err }, errEmptyGroupKey},
		{"Group.Delete", func() error { _, err := client.Group.Delete(ctx, ""); return err }, errEmptyGroupKey},
		{"Group.Feeds", func() error { _, _, err := client.Group.Feeds(ctx, ""); return err }, errEmptyGroupKey},
		{"Group.CreateFeed", func() error { _, _, err := client.Group.CreateFeed(ctx, "", &Feed{}); return err }, errEmptyGroupKey},
		{"Group.AddFeed", func() error { _, _, err := client.Group.AddFeed(ctx, "", "f"); return err }, errEmptyGroupKey},
		{"Group.AddFeed feed", func() error { _, _, err := client.Group.AddFeed(ctx, "g", ""); return err }, errEmptyFeedKey},
		{"Group.RemoveFeed", func() error { _, _, err := client.Group.RemoveFeed(ctx, "", "f"); return err }, errEmptyGroupKey},
		{"Group.RemoveFeed feed", func() error { _, _, err := client.Group.RemoveFeed(ctx, "g", ""); return err }, errEmptyFeedKey},
		{"Group.LastValues", func() error { _, _, err := client.Group.LastValues(ctx, ""); return err }, errEmptyGroupKey},
		{"Data.Get ID", func() error { _, _, err := client.Data.Get(ctx, "f", ""); return err }, errEmptyID},
		{"Data.Update ID", func() error { _, _, err := client.Data.Update(ctx, "f", "", &Data{}); return err }, errEmptyID},
		{"Data.Delete ID", func() error { _, err := client.Data.Delete(ctx, "f", ""); return err }, errEmptyID},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.call(), tt.name)
	}
}
