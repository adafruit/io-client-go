package adafruitio

import (
	"encoding/json"
	"fmt"
)

type DataService struct {
	client *Client
}

// DataPoint are what we create in the data service
type DataPoint struct {
	ID           int         `json:"id,omitempty"`
	Value        json.Number `json:"value,omitempty"` // number, string, ?
	Position     string      `json:"position,omitempty"`
	FeedID       int         `json:"feed_id,omitempty"`
	GroupID      int         `json:"group_id,omitempty"`
	Expiration   string      `json:"expiration,omitempty"`
	Latitude     float64     `json:"lat,omitempty"`
	Longitude    float64     `json:"lon,omitempty"`
	Elevation    float64     `json:"ele,omitempty"`
	CompletedAt  string      `json:"completed_at,omitempty"`
	CreatedAt    string      `json:"created_at,omitempty"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
	CreatedEpoch float64     `json:"created_epoch,omitempty"`
}

// GET /feeds/{feed_id}/data
//
// Get all Data for an existing Fees.
func (s *DataService) All() ([]*DataPoint, *Response, error) {
	path, ferr := s.client.Feed.Path("/data")
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	datas := make([]*DataPoint, 0)
	resp, err := s.client.Do(req, &datas)
	if err != nil {
		return nil, nil, err
	}

	return datas, resp, nil
}

// Get(id int)
func (s *DataService) Get(id int) (*DataPoint, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", id))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var data DataPoint
	resp, err := s.client.Do(req, &data)
	if err != nil {
		return nil, nil, err
	}

	return &data, resp, nil
}

// PATCH /feeds/{feed_id}/data/{data_id}
//
// Update takes an ID and a Data record, updates the record idendified by ID,
// and returns an new, updated record instance or an error.
func (s *DataService) Update(id interface{}, data *DataPoint) (*DataPoint, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", id))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("PATCH", path, data)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedData DataPoint
	resp, err := s.client.Do(req, &updatedData)
	if err != nil {
		return nil, nil, err
	}

	return &updatedData, resp, nil
}

// Delete(id int)
//
// Delete the Data identified by the given ID.
func (s *DataService) Delete(id int) (*Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", id))
	if ferr != nil {
		return nil, ferr
	}

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

// private method for handling the Next, Prev, and Last commands
func (s *DataService) retrieve(command string) (*DataPoint, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", command))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var data DataPoint
	resp, err := s.client.Do(req, &data)
	if err != nil {
		return nil, nil, err
	}

	return &data, resp, nil
}

// Get the next Data in the queue.
func (s *DataService) Next() (*DataPoint, *Response, error) {
	return s.retrieve("next")
}

// Get the previous Data in the queue.
func (s *DataService) Prev() (*DataPoint, *Response, error) {
	return s.retrieve("prev")
}

// Get the last Data in the queue.
func (s *DataService) Last() (*DataPoint, *Response, error) {
	return s.retrieve("last")
}

// POST /feeds/{feed_id}/data
//
// Create new Data on an existing Feed
func (s *DataService) Create(dp *DataPoint) (*DataPoint, *Response, error) {
	path, ferr := s.client.Feed.Path("/data")
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates a new datapoint
	point := &DataPoint{}
	resp, err := s.client.Do(req, point)
	if err != nil {
		return nil, nil, err
	}

	return point, resp, nil
}

// POST /feeds/{feed_id}/send
//
// Create new Data point on the CurrentFeed, also create the Feed if necessary.
func (s *DataService) Send(dp *DataPoint) (*DataPoint, *Response, error) {
	path, ferr := s.client.Feed.Path("/data/send")
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	point := &DataPoint{}
	resp, err := s.client.Do(req, point)
	if err != nil {
		return nil, nil, err
	}

	return point, resp, nil
}
