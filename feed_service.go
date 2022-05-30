package adafruitio

import (
	"fmt"
	"path"
)

type FeedService struct {
	// CurrentFeed is the Feed used for all Data access.
	CurrentFeed *Feed

	client *Client
}

// Path generates a Feed-specific path with the given suffix.
func (s *FeedService) Path(suffix string) (string, error) {
	ferr := s.client.checkFeed()
	if ferr != nil {
		return "", ferr
	}
	return path.Join(fmt.Sprintf("feeds/%v", s.CurrentFeed.Key), suffix), nil
}

type Owner struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type Feed struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Key           string `json:"key,omitempty"`
	Username      string `json:"username,omitempty"`
	Owner         *Owner `json:"owner,omitempty"`
	Description   string `json:"description,omitempty"`
	UnitType      string `json:"unit_type,omitempty"`
	UnitSymbol    string `json:"unit_symbol,omitempty"`
	History       bool   `json:"history,omitempty"`
	Visibility    string `json:"visibility,omitempty"`
	License       string `json:"license,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
	LastValue     string `json:"last_value,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusNotify  bool   `json:"status_notify,omitempty"`
	StatusTimeout int    `json:"status_timeout,omitempty"`
	Shared        bool   `json:"is_shared,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

// All lists all available feeds.
func (s *FeedService) All() ([]*Feed, *Response, error) {
	path := "feeds"

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	feeds := make([]*Feed, 0)
	resp, err := s.client.Do(req, &feeds)
	if err != nil {
		return nil, resp, err
	}

	return feeds, resp, nil
}

// Get returns the Feed record identified by the given parameter. Parameter can
// be the Feed's Name, Key, or ID.
func (s *FeedService) Get(key string) (*Feed, *Response, error) {
	path := fmt.Sprintf("feeds/%s", key)

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var feed Feed
	resp, err := s.client.Do(req, &feed)
	if err != nil {
		return nil, resp, err
	}

	return &feed, resp, nil
}

// Create takes a Feed record, creates it, and returns the updated record or an error.
func (s *FeedService) Create(feed *Feed) (*Feed, *Response, error) {
	path := "feeds"

	req, rerr := s.client.NewRequest("POST", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	resp, err := s.client.Do(req, feed)
	if err != nil {
		return nil, resp, err
	}

	return feed, resp, nil
}

// Update takes an ID and a Feed record, updates it, and returns an updated
// record instance or an error.
//
// Only the Feed Name and Description can be modified.
func (s *FeedService) Update(key string, feed *Feed) (*Feed, *Response, error) {
	path := fmt.Sprintf("feeds/%s", key)

	req, rerr := s.client.NewRequest("PATCH", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedFeed Feed
	resp, err := s.client.Do(req, &updatedFeed)
	if err != nil {
		return nil, resp, err
	}

	return &updatedFeed, resp, nil
}

// Delete the Feed identified by the given ID.
func (s *FeedService) Delete(key string) (*Response, error) {
	path := fmt.Sprintf("feeds/%s", key)

	req, rerr := s.client.NewRequest("DELETE", path, nil)
	if rerr != nil {
		return nil, rerr
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
