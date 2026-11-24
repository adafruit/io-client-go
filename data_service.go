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
//
// Results are sorted newest first and paginated. To get every page, set
// Before to the NextBefore of the previous page's Response until it is empty.
type DataFilter struct {
	// StartTime returns only Data created after this time. If StartTime is
	// set without EndTime, the API returns the oldest records after
	// StartTime, so set EndTime too when paging through a time range.
	StartTime time.Time `url:"start_time,omitempty"`

	// EndTime returns only Data created before this time.
	EndTime time.Time `url:"end_time,omitempty"`

	// Limit is the number of records per page. The API's default is 100 and
	// its maximum is 1000.
	Limit int `url:"limit,omitempty"`

	// Include lists the Data fields to return, for example
	// []string{"value", "created_at"}. Fields that aren't included are nil.
	// The API accepts value, lat, lon, ele, id, and created_at.
	Include []string `url:"include,omitempty,comma"`

	// Before is the cursor for the page to get. Set it to the NextBefore of
	// the previous page's Response.
	Before string `url:"before,omitempty"`
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

// CreateBatch adds several Data values to the Feed identified by feedKey in a
// single request, and returns the created records.
func (s *DataService) CreateBatch(ctx context.Context, feedKey string, data []*Data) ([]*Data, *Response, error) {
	path, perr := dataPath(feedKey, "batch")
	if perr != nil {
		return nil, nil, perr
	}

	body := struct {
		Data []*Data `json:"data"`
	}{data}

	req, rerr := s.client.NewRequest("POST", path, body)
	if rerr != nil {
		return nil, nil, rerr
	}

	created := make([]*Data, 0)
	resp, err := s.client.Do(ctx, req, &created)
	if err != nil {
		return nil, resp, err
	}

	return created, resp, nil
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

// LastValue returns the most recent Data in the Feed identified by feedKey.
// Unlike Last, it works with read-only access, so it can be used on public
// feeds. If the feed has no data, the API returns a 404 error.
func (s *DataService) LastValue(ctx context.Context, feedKey string) (*Data, *Response, error) {
	return s.retrieve(ctx, feedKey, "last_value?type=json")
}

// ChartFilter sets the time range and aggregation for DataService.Chart.
// Zero values are left out of the request.
type ChartFilter struct {
	// StartTime and EndTime set the time range of the chart. If both are
	// set, Hours is ignored.
	StartTime time.Time `url:"start_time,omitempty"`
	EndTime   time.Time `url:"end_time,omitempty"`

	// Hours sets the time range to the given number of hours before now.
	Hours int `url:"hours,omitempty"`

	// Resolution is the size of each aggregation slice in minutes: 1, 5, 10,
	// 30, 60, 120, 240, 480, or 960. If it is 0, the API picks one based on
	// the time range.
	Resolution int `url:"resolution,omitempty"`

	// Field is the aggregate to return for each slice: avg, sum, val, min,
	// max, or val_count. The API's default is avg.
	Field string `url:"field,omitempty"`

	// Raw returns the raw data points instead of aggregates. It can't be
	// used with Resolution or Field.
	Raw bool `url:"raw,omitempty"`
}

// Chart is chart data for a Feed, as returned by DataService.Chart.
type Chart struct {
	// Feed identifies the charted feed by ID, Key, and Name.
	Feed *Feed `json:"feed,omitempty"`

	// Parameters are the parameters the API used to build the chart,
	// including any it picked itself, such as Resolution.
	Parameters *ChartParameters `json:"parameters,omitempty"`

	// Columns names the values in each row of Data, for example
	// ["date", "value"].
	Columns []string `json:"columns,omitempty"`

	// Data has one row per point, with values in the order of Columns. Dates
	// are RFC 3339 strings.
	Data [][]string `json:"data,omitempty"`

	// Storage reports whether the chart was built from raw ("raw") or
	// aggregated ("agg") data.
	Storage *string `json:"storage,omitempty"`
}

// ChartParameters are the parameters used to build a Chart.
type ChartParameters struct {
	StartTime  *Timestamp `json:"start_time,omitempty"`
	EndTime    *Timestamp `json:"end_time,omitempty"`
	Hours      *int       `json:"hours,omitempty"`
	Resolution *int       `json:"resolution,omitempty"`
	Field      *string    `json:"field,omitempty"`
}

// Chart returns chart data for the Feed identified by feedKey: the feed's
// values over a time range, aggregated into time slices. opt may be nil.
func (s *DataService) Chart(ctx context.Context, feedKey string, opt *ChartFilter) (*Chart, *Response, error) {
	path, perr := dataPath(feedKey, "chart")
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

	var chart Chart
	resp, err := s.client.Do(ctx, req, &chart)
	if err != nil {
		return nil, resp, err
	}

	return &chart, resp, nil
}
