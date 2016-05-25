package adafruitio

import (
	"encoding/json"
	"path"
)

type FeedService struct {
	client *Client

	// the name or ID of the feed
	Name string
}

type Feed struct {
	ID          json.Number `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Key         string      `json:"key,omitempty"`
	Description string      `json:"description,omitempty"`
	UnitType    string      `json:"unit_type,omitempty"`
	UnitSymbol  string      `json:"unit_symbol,omitempty"`
	History     bool        `json:"history,omitempty"`
	Visibility  string      `json:"visibility,omitempty"`
	License     string      `json:"license,omitempty"`
	Enabled     bool        `json:"enabled,omitempty"`
	LastValue   string      `json:"last_value,omitempty"`
	Status      string      `json:"status,omitempty"`
	GroupID     json.Number `json:"group_id,omitempty"`
	CreatedAt   string      `json:"created_at,omitempty"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
}

func (s *FeedService) All() ([]Feed, *Response, error) {
	path := path.Join("api", "v1", "feeds")

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	feeds := make([]Feed, 0)
	resp, err := s.client.Do(req, &feeds)
	if err != nil {
		return nil, nil, err
	}

	return feeds, resp, nil
}
