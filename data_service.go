package adafruitio

import (
	"context"
	"errors"
	"path"
	"time"
)

// Data is a single value in a Feed. See Feed for how pointer fields work.
type Data struct {
	ID           *string    `json:"id,omitempty"`
	Value        *string    `json:"value,omitempty"`
	FeedID       *int64     `json:"feed_id,omitempty"`
	FeedKey      *string    `json:"feed_key,omitempty"`
	GroupID      *int64     `json:"group_id,omitempty"`
	Expiration   *Timestamp `json:"expiration,omitempty"`
	Latitude     *float64   `json:"lat,omitempty"`
	Longitude    *float64   `json:"lon,omitempty"`
	Elevation    *float64   `json:"ele,omitempty"`
	CompletedAt  *Timestamp `json:"completed_at,omitempty"`
	CreatedAt    *Timestamp `json:"created_at,omitempty"`
	UpdatedAt    *Timestamp `json:"updated_at,omitempty"`
	CreatedEpoch *float64   `json:"created_epoch,omitempty"`
}

// DataFilter filters the Data returned by DataService.All and Search. Zero
// values are left out of the request.
type DataFilter struct {
	// StartTime returns only Data created after this time.
	StartTime time.Time `url:"start_time,omitempty"`
	// EndTime returns only Data created before this time.
	EndTime time.Time `url:"end_time,omitempty"`
}

type DataService struct {
	client *Client
}

var errEmptyFeedKey = errors.New("feed key must not be empty")

// dataPath builds the path for a feed's data endpoints, e.g.
// dataPath("temperature", "last") returns "feeds/temperature/data/last".
func dataPath(feedKey, suffix string) (string, error) {
	if feedKey == "" {
		return "", errEmptyFeedKey
	}
	return path.Join("feeds", feedKey, "data", suffix), nil
}

// All returns Data for the Feed identified by feedKey, optionally filtered by
// opt. opt may be nil.
func (s *DataService) All(ctx context.Context, feedKey string, opt *DataFilter) ([]*Data, *Response, error) {
	path, perr := dataPath(feedKey, "")
	if perr != nil {
		return nil, nil, perr
	}

	path, oerr := addOptions(path, opt)
	if oerr != nil {
		return nil, nil, oerr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	datas := make([]*Data, 0)
	resp, err := s.client.Do(ctx, req, &datas)
	if err != nil {
		return nil, resp, err
	}

	return datas, resp, nil
}

// Search returns Data for the Feed identified by feedKey that matches the
// given filter. It is equivalent to All.
func (s *DataService) Search(ctx context.Context, feedKey string, filter *DataFilter) ([]*Data, *Response, error) {
	return s.All(ctx, feedKey, filter)
}

// Get returns the Data record identified by id from the Feed identified by
// feedKey.
func (s *DataService) Get(ctx context.Context, feedKey, id string) (*Data, *Response, error) {
	return s.retrieve(ctx, feedKey, id)
}

// Update updates the Data record identified by id in the Feed identified by
// feedKey, and returns the updated record.
func (s *DataService) Update(ctx context.Context, feedKey, id string, data *Data) (*Data, *Response, error) {
	path, perr := dataPath(feedKey, id)
	if perr != nil {
		return nil, nil, perr
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

// Delete deletes the Data record identified by id from the Feed identified by
// feedKey.
func (s *DataService) Delete(ctx context.Context, feedKey, id string) (*Response, error) {
	path, perr := dataPath(feedKey, id)
	if perr != nil {
		return nil, perr
	}

	req, rerr := s.client.NewRequest("DELETE", path, nil)
	if rerr != nil {
		return nil, rerr
	}

	return s.client.Do(ctx, req, nil)
}

// retrieve GETs a single Data record at feeds/{feedKey}/data/{suffix}, where
// suffix is a record ID or one of the next, previous, first, and last
// commands.
func (s *DataService) retrieve(ctx context.Context, feedKey, suffix string) (*Data, *Response, error) {
	path, perr := dataPath(feedKey, suffix)
	if perr != nil {
		return nil, nil, perr
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

// Next returns the next Data in the Feed identified by feedKey.
func (s *DataService) Next(ctx context.Context, feedKey string) (*Data, *Response, error) {
	return s.retrieve(ctx, feedKey, "next")
}

// Prev returns the previous Data in the Feed identified by feedKey.
func (s *DataService) Prev(ctx context.Context, feedKey string) (*Data, *Response, error) {
	return s.retrieve(ctx, feedKey, "previous")
}

// First returns the first Data in the Feed identified by feedKey.
func (s *DataService) First(ctx context.Context, feedKey string) (*Data, *Response, error) {
	return s.retrieve(ctx, feedKey, "first")
}

// Last returns the last Data in the Feed identified by feedKey.
func (s *DataService) Last(ctx context.Context, feedKey string) (*Data, *Response, error) {
	return s.retrieve(ctx, feedKey, "last")
}

// Create adds a new Data value to the Feed identified by feedKey, and returns
// the created record.
func (s *DataService) Create(ctx context.Context, feedKey string, dp *Data) (*Data, *Response, error) {
	path, perr := dataPath(feedKey, "")
	if perr != nil {
		return nil, nil, perr
	}

	req, rerr := s.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	point := &Data{}
	resp, err := s.client.Do(ctx, req, point)
	if err != nil {
		return nil, resp, err
	}

	return point, resp, nil
}
