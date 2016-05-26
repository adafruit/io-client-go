package adafruitio

import (
	"fmt"
	"path"
)

type FeedService struct {
	CurrentFeed *Feed

	client *Client
}

func (s *FeedService) Path(part string) (string, error) {
	ferr := s.client.checkFeed()
	if ferr != nil {
		return "", ferr
	}
	return path.Join(fmt.Sprintf("api/v1/feeds/%v", s.CurrentFeed.Key), part), nil
}

type Feed struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
	UnitType    string `json:"unit_type,omitempty"`
	UnitSymbol  string `json:"unit_symbol,omitempty"`
	History     bool   `json:"history,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	License     string `json:"license,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	LastValue   string `json:"last_value,omitempty"`
	Status      string `json:"status,omitempty"`
	GroupID     int    `json:"group_id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// All lists all available feeds.
func (s *FeedService) All() ([]*Feed, *Response, error) {
	path := "api/v1/feeds"

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	feeds := make([]*Feed, 0)
	resp, err := s.client.Do(req, &feeds)
	if err != nil {
		return nil, nil, err
	}

	return feeds, resp, nil
}

// Get the Feed record identified by the given ID
func (s *FeedService) Get(id interface{}) (*Feed, *Response, error) {
	path := fmt.Sprintf("api/v1/feeds/%v", id)

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var feed Feed
	resp, err := s.client.Do(req, &feed)
	if err != nil {
		return nil, nil, err
	}

	return &feed, resp, nil
}

// Create takes a Feed record, creates it, and returns the updated record or an error.
func (s *FeedService) Create(feed *Feed) (*Feed, *Response, error) {
	path := "api/v1/feeds"

	req, rerr := s.client.NewRequest("POST", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	resp, err := s.client.Do(req, feed)
	if err != nil {
		return nil, nil, err
	}

	return feed, resp, nil
}

// Update takes an ID and a Feed record, updates it, and returns an updated
// record instance or an error.
func (s *FeedService) Update(id interface{}, feed *Feed) (*Feed, *Response, error) {
	path := fmt.Sprintf("api/v1/feeds/%v", id)

	req, rerr := s.client.NewRequest("PATCH", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedFeed Feed
	resp, err := s.client.Do(req, &updatedFeed)
	if err != nil {
		return nil, nil, err
	}

	return &updatedFeed, resp, nil
}

// Delete the Feed identified by the given ID.
func (s *FeedService) Delete(id interface{}) (*Response, error) {
	path := fmt.Sprintf("api/v1/feeds/%v", id)

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
