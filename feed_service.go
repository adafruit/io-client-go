package adafruitio

import (
	"context"
	"fmt"
)

type FeedService struct {
	client *Client
}

// Owner identifies the user who owns a Feed or Group.
type Owner struct {
	ID       *int64  `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
}

// Feed is an Adafruit IO feed.
//
// As in go-github, every field is a pointer. A nil field is left out of create
// and update requests, so the server keeps its current value, while a non-nil
// field is always sent, even if it is empty or false. Use Ptr to set fields,
// for example &Feed{Description: Ptr(""), History: Ptr(false)}, and the
// generated Get methods, such as GetName, to read them without nil checks.
type Feed struct {
	ID            *int64  `json:"id,omitempty"`
	Name          *string `json:"name,omitempty"`
	Key           *string `json:"key,omitempty"`
	Username      *string `json:"username,omitempty"`
	Owner         *Owner  `json:"owner,omitempty"`
	Description   *string `json:"description,omitempty"`
	UnitType      *string `json:"unit_type,omitempty"`
	UnitSymbol    *string `json:"unit_symbol,omitempty"`
	History       *bool   `json:"history,omitempty"`
	Visibility    *string `json:"visibility,omitempty"`
	License       *string `json:"license,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
	LastValue     *string `json:"last_value,omitempty"`
	Status        *string `json:"status,omitempty"`
	StatusNotify  *bool   `json:"status_notify,omitempty"`
	StatusTimeout *int    `json:"status_timeout,omitempty"`
	CreatedAt     *string `json:"created_at,omitempty"`
	UpdatedAt     *string `json:"updated_at,omitempty"`

	// Writable reports whether the current user can write to the feed. It is
	// read-only, and only returned when getting a single feed.
	Writable *bool `json:"writable,omitempty"`

	// Group is the feed's primary group, and Groups lists every group the
	// feed belongs to. Both are read-only.
	Group  *Group   `json:"group,omitempty"`
	Groups []*Group `json:"groups,omitempty"`
}

// All lists all available feeds.
func (s *FeedService) All(ctx context.Context) ([]*Feed, *Response, error) {
	path := "feeds"

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	feeds := make([]*Feed, 0)
	resp, err := s.client.Do(ctx, req, &feeds)
	if err != nil {
		return nil, resp, err
	}

	return feeds, resp, nil
}

// Get returns the Feed record identified by the given parameter. Parameter can
// be the Feed's Name, Key, or ID.
func (s *FeedService) Get(ctx context.Context, key string) (*Feed, *Response, error) {
	path := fmt.Sprintf("feeds/%s", key)

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var feed Feed
	resp, err := s.client.Do(ctx, req, &feed)
	if err != nil {
		return nil, resp, err
	}

	return &feed, resp, nil
}

// Create takes a Feed record, creates it, and returns the updated record or an error.
func (s *FeedService) Create(ctx context.Context, feed *Feed) (*Feed, *Response, error) {
	path := "feeds"

	req, rerr := s.client.NewRequest("POST", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	resp, err := s.client.Do(ctx, req, feed)
	if err != nil {
		return nil, resp, err
	}

	return feed, resp, nil
}

// Update takes an ID and a Feed record, updates it, and returns an updated
// record instance or an error.
//
// Only the Feed Name and Description can be modified.
func (s *FeedService) Update(ctx context.Context, key string, feed *Feed) (*Feed, *Response, error) {
	path := fmt.Sprintf("feeds/%s", key)

	req, rerr := s.client.NewRequest("PATCH", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedFeed Feed
	resp, err := s.client.Do(ctx, req, &updatedFeed)
	if err != nil {
		return nil, resp, err
	}

	return &updatedFeed, resp, nil
}

// Delete the Feed identified by the given ID.
func (s *FeedService) Delete(ctx context.Context, key string) (*Response, error) {
	path := fmt.Sprintf("feeds/%s", key)

	req, rerr := s.client.NewRequest("DELETE", path, nil)
	if rerr != nil {
		return nil, rerr
	}

	return s.client.Do(ctx, req, nil)
}
