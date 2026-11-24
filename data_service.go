package adafruitio

import (
	"context"
	"fmt"
)

// Data are the values contained by a Feed.
type Data struct {
	ID           string  `json:"id,omitempty"`
	Value        string  `json:"value,omitempty"`
	FeedID       int     `json:"feed_id,omitempty"`
	FeedKey      string  `json:"feed_key,omitempty"`
	GroupID      int     `json:"group_id,omitempty"`
	Expiration   string  `json:"expiration,omitempty"`
	Latitude     float64 `json:"lat,omitempty"`
	Longitude    float64 `json:"lon,omitempty"`
	Elevation    float64 `json:"ele,omitempty"`
	CompletedAt  string  `json:"completed_at,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
	CreatedEpoch float64 `json:"created_epoch,omitempty"`
}

type DataFilter struct {
	StartTime string `url:"start_time,omitempty"`
	EndTime   string `url:"end_time,omitempty"`
}

type DataService struct {
	client *Client
}

// All returns all Data for the currently selected Feed. See Client.SetFeed()
// for details on selecting a Feed.
func (s *DataService) All(ctx context.Context, opt *DataFilter) ([]*Data, *Response, error) {
	path, ferr := s.client.Feed.Path("/data")
	if ferr != nil {
		return nil, nil, ferr
	}

	path, oerr := addOptions(path, opt)
	if oerr != nil {
		return nil, nil, oerr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	datas := make([]*Data, 0)
	resp, err := s.client.Do(ctx, req, &datas)
	if err != nil {
		return nil, resp, err
	}

	return datas, resp, nil
}

// Search returns Data for the currently selected Feed that matches the given
// filter. It is equivalent to All.
func (s *DataService) Search(ctx context.Context, filter *DataFilter) ([]*Data, *Response, error) {
	return s.All(ctx, filter)
}

// Get returns a single Data element, identified by the given ID parameter.
func (s *DataService) Get(ctx context.Context, id string) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%s", id))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var data Data
	resp, err := s.client.Do(ctx, req, &data)
	if err != nil {
		return nil, resp, err
	}

	return &data, resp, nil
}

// Update takes an ID and a Data record, updates the record idendified by ID,
// and returns a new, updated Data instance.
func (s *DataService) Update(ctx context.Context, id string, data *Data) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%s", id))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("PATCH", path, data)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedData Data
	resp, err := s.client.Do(ctx, req, &updatedData)
	if err != nil {
		return nil, resp, err
	}

	return &updatedData, resp, nil
}

// Delete the Data identified by the given ID.
func (s *DataService) Delete(ctx context.Context, id string) (*Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%s", id))
	if ferr != nil {
		return nil, ferr
	}

	req, rerr := s.client.NewRequest("DELETE", path, nil)
	if rerr != nil {
		return nil, rerr
	}

	return s.client.Do(ctx, req, nil)
}

// private method for handling the Next, Prev, and Last commands
func (s *DataService) retrieve(ctx context.Context, command string) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", command))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var data Data
	resp, err := s.client.Do(ctx, req, &data)
	if err != nil {
		return nil, resp, err
	}

	return &data, resp, nil
}

// Next returns the next Data in the stream.
func (s *DataService) Next(ctx context.Context) (*Data, *Response, error) {
	return s.retrieve(ctx, "next")
}

// Prev returns the previous Data in the stream.
func (s *DataService) Prev(ctx context.Context) (*Data, *Response, error) {
	return s.retrieve(ctx, "previous")
}

// First returns the first Data in the stream.
func (s *DataService) First(ctx context.Context) (*Data, *Response, error) {
	return s.retrieve(ctx, "first")
}

// Last returns the last Data in the stream.
func (s *DataService) Last(ctx context.Context) (*Data, *Response, error) {
	return s.retrieve(ctx, "last")
}

// Create adds a new Data value to an existing Feed.
func (s *DataService) Create(ctx context.Context, dp *Data) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path("/data")
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates a new datapoint
	point := &Data{}
	resp, err := s.client.Do(ctx, req, point)
	if err != nil {
		return nil, resp, err
	}

	return point, resp, nil
}
