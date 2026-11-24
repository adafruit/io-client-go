package adafruitio

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataUnmarshalLocation(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		lat, lon, ele *float64
	}{
		// create and batch return numbers
		{"numbers", `{"lat": 40.7, "lon": -74.0, "ele": 10}`, Ptr(40.7), Ptr(-74.0), Ptr(10.0)},
		// list, get, first, last, previous, and last_value return strings
		{"strings", `{"lat": "40.7", "lon": "-74.0", "ele": "10.0"}`, Ptr(40.7), Ptr(-74.0), Ptr(10.0)},
		{"nulls", `{"lat": null, "lon": null, "ele": null}`, nil, nil, nil},
		{"empty strings", `{"lat": "", "lon": "", "ele": ""}`, nil, nil, nil},
		{"missing", `{}`, nil, nil, nil},
	}

	for _, tt := range tests {
		var d Data
		err := json.Unmarshal([]byte(tt.json), &d)

		assert.NoError(t, err, tt.name)
		assert.Equal(t, tt.lat, d.Latitude, tt.name)
		assert.Equal(t, tt.lon, d.Longitude, tt.name)
		assert.Equal(t, tt.ele, d.Elevation, tt.name)
	}
}

func TestDataUnmarshalInvalidLocation(t *testing.T) {
	var d Data
	assert.Error(t, json.Unmarshal([]byte(`{"lat": "north"}`), &d))
	assert.Error(t, json.Unmarshal([]byte(`{"lat": true}`), &d))
}

// Decoding a full record from the list endpoint, as returned by the API.
func TestDataUnmarshalListRecord(t *testing.T) {
	var data []*Data
	err := json.Unmarshal([]byte(`[{
		"id": "0GAJ9TJ8004DMKAYN8QBF1S65G",
		"value": "21.5",
		"feed_id": 32,
		"feed_key": "probe-loc",
		"created_at": "2026-11-24T18:14:24Z",
		"created_epoch": 1795544064,
		"expiration": "2027-01-23T18:14:24Z",
		"lat": "40.7",
		"lon": "-74.0",
		"ele": "10.0",
		"location": {"type": "Feature", "geometry": {"type": "Point", "coordinates": ["-74.0", "40.7", "10.0"]}}
	}]`), &data)

	assert.NoError(t, err)
	if assert.Len(t, data, 1) {
		d := data[0]
		assert.Equal(t, "21.5", d.GetValue())
		assert.Equal(t, int64(32), d.GetFeedID())
		assert.Equal(t, 40.7, d.GetLatitude())
		assert.Equal(t, -74.0, d.GetLongitude())
		assert.Equal(t, 10.0, d.GetElevation())
		assert.False(t, d.GetCreatedAt().IsZero())
	}
}

// Fields that aren't in the JSON are left unchanged, like encoding/json does.
func TestDataUnmarshalKeepsMissingFields(t *testing.T) {
	d := Data{Latitude: Ptr(1.5), Value: Ptr("old")}

	err := json.Unmarshal([]byte(`{"value": "new"}`), &d)

	assert.NoError(t, err)
	assert.Equal(t, "new", d.GetValue())
	assert.Equal(t, Ptr(1.5), d.Latitude)
}

// Encoding is unaffected: lat, lon, and ele are sent as numbers.
func TestDataMarshalLocation(t *testing.T) {
	b, err := json.Marshal(&Data{Value: Ptr("1"), Latitude: Ptr(40.7), Longitude: Ptr(-74.0)})

	assert.NoError(t, err)
	assert.JSONEq(t, `{"value": "1", "lat": 40.7, "lon": -74.0}`, string(b))
}
