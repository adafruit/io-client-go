package adafruitio

import (
	"encoding/json"
	"path"
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

// POST /feeds/{feed_id}/data
//
// Create new Data on an existing Feed
func (d *DataService) Create(dp *DataPoint) (*DataPoint, *Response, error) {
	// feed name must be set before Data interface can be called
	ferr := d.client.checkFeed()
	if ferr != nil {
		return nil, nil, ferr
	}

	path := path.Join("api", "v1", "feeds", d.client.Feed.Name, "data")

	req, rerr := d.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates a new datapoint
	point := &DataPoint{}
	resp, err := d.client.Do(req, point)
	if err != nil {
		return nil, nil, err
	}

	return point, resp, nil
}

// POST /feeds/{feed_id}/send
//
// Create new Data and Feed.
func (d *DataService) Send(dp *DataPoint) (*DataPoint, *Response, error) {
	ferr := d.client.checkFeed()
	if ferr != nil {
		return nil, nil, ferr
	}

	path := path.Join("api", "v1", "feeds", d.client.Feed.Name, "data", "send")

	req, rerr := d.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	point := &DataPoint{}
	resp, err := d.client.Do(req, point)
	if err != nil {
		return nil, nil, err
	}

	return point, resp, nil
}
