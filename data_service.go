package adafruitio

import "fmt"

// Data are the values contained by a Feed.
type Data struct {
	ID           int     `json:"id,omitempty"`
	Value        string  `json:"value,omitempty"`
	Position     string  `json:"position,omitempty"`
	FeedID       int     `json:"feed_id,omitempty"`
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

type DataService struct {
	client *Client
}

// All returns all Data for the currently selected Feed. See Client.SetFeed()
// for details on selecting a Feed.
func (s *DataService) All() ([]*Data, *Response, error) {
	path, ferr := s.client.Feed.Path("/data")
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	datas := make([]*Data, 0)
	resp, err := s.client.Do(req, &datas)
	if err != nil {
		return nil, resp, err
	}

	return datas, resp, nil
}

// Get returns a single Data element, identified by the given ID parameter.
func (s *DataService) Get(id int) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", id))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var data Data
	resp, err := s.client.Do(req, &data)
	if err != nil {
		return nil, resp, err
	}

	return &data, resp, nil
}

// Update takes an ID and a Data record, updates the record idendified by ID,
// and returns a new, updated Data instance.
func (s *DataService) Update(id interface{}, data *Data) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", id))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("PATCH", path, data)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedData Data
	resp, err := s.client.Do(req, &updatedData)
	if err != nil {
		return nil, resp, err
	}

	return &updatedData, resp, nil
}

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
func (s *DataService) retrieve(command string) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path(fmt.Sprintf("/data/%v", command))
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var data Data
	resp, err := s.client.Do(req, &data)
	if err != nil {
		return nil, resp, err
	}

	return &data, resp, nil
}

// Next returns the next Data in the stream.
func (s *DataService) Next() (*Data, *Response, error) {
	return s.retrieve("next")
}

// Prev returns the previous Data in the stream.
func (s *DataService) Prev() (*Data, *Response, error) {
	return s.retrieve("prev")
}

// Last returns the last Data in the stream.
func (s *DataService) Last() (*Data, *Response, error) {
	return s.retrieve("last")
}

// Create adds a new Data value to an existing Feed.
func (s *DataService) Create(dp *Data) (*Data, *Response, error) {
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
	resp, err := s.client.Do(req, point)
	if err != nil {
		return nil, resp, err
	}

	return point, resp, nil
}

// Send adds a new Data value to an existing Feed, or will create the Feed if
// it doesn't already exist.
func (s *DataService) Send(dp *Data) (*Data, *Response, error) {
	path, ferr := s.client.Feed.Path("/data/send")
	if ferr != nil {
		return nil, nil, ferr
	}

	req, rerr := s.client.NewRequest("POST", path, dp)
	if rerr != nil {
		return nil, nil, rerr
	}

	point := &Data{}
	resp, err := s.client.Do(req, point)
	if err != nil {
		return nil, resp, err
	}

	return point, resp, nil
}
