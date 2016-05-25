package adafruitio

import "fmt"

type FeedService struct {
	client *Client

	// the name or ID of the feed
	Name string
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

func (s *FeedService) Get(id interface{}) (*Feed, *Response, error) {
	path := fmt.Sprintf("api/v1/feeds/%v", id)

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	var feed Feed
	resp, err := s.client.Do(req, &feed)
	if err != nil {
		return nil, nil, err
	}

	return &feed, resp, nil
}

func (s *FeedService) Create(id int) (*Feed, *Response, error) {
	return nil, nil, nil
}

func (s *FeedService) Update(id int) (*Feed, *Response, error) {
	return nil, nil, nil
}

func (s *FeedService) Delete(id int) (*Feed, *Response, error) {
	return nil, nil, nil
}
