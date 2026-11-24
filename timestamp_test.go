package adafruitio

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    time.Time
		wantErr bool
	}{
		{"RFC 3339", `"2026-11-24T16:05:19Z"`, time.Date(2026, 11, 24, 16, 5, 19, 0, time.UTC), false},
		{"fractional seconds", `"2026-01-01T00:00:00.123Z"`, time.Date(2026, 1, 1, 0, 0, 0, 123e6, time.UTC), false},
		{"offset", `"2026-11-24T10:05:19-06:00"`, time.Date(2026, 11, 24, 16, 5, 19, 0, time.UTC), false},
		{"Unix seconds", `1795536337`, time.Unix(1795536337, 0), false},
		{"fractional Unix seconds", `1795536337.5`, time.Unix(1795536337, 5e8), false},
		{"quoted Unix seconds", `"1800725486.0"`, time.Unix(1800725486, 0), false},
		{"empty string", `""`, time.Time{}, false},
		{"null", `null`, time.Time{}, false},
		{"invalid string", `"yesterday"`, time.Time{}, true},
		{"invalid type", `true`, time.Time{}, true},
	}

	for _, tt := range tests {
		var ts Timestamp
		err := json.Unmarshal([]byte(tt.json), &ts)

		if tt.wantErr {
			assert.Error(t, err, tt.name)
			continue
		}
		assert.NoError(t, err, tt.name)
		assert.True(t, ts.Time.Equal(tt.want), "%s: got %v, want %v", tt.name, ts.Time, tt.want)
	}
}

func TestTimestampMarshalJSON(t *testing.T) {
	ts := Timestamp{time.Date(2026, 11, 24, 16, 5, 19, 0, time.UTC)}

	b, err := json.Marshal(ts)

	assert.NoError(t, err)
	assert.Equal(t, `"2026-11-24T16:05:19Z"`, string(b))

	// and it round-trips
	var decoded Timestamp
	assert.NoError(t, json.Unmarshal(b, &decoded))
	assert.True(t, ts.Equal(decoded))
}

func TestTimestampGetTime(t *testing.T) {
	var nilTimestamp *Timestamp
	assert.Nil(t, nilTimestamp.GetTime())

	ts := &Timestamp{time.Date(2026, 11, 24, 16, 5, 19, 0, time.UTC)}
	assert.Equal(t, ts.Time, *ts.GetTime())
}

func TestDataDecodesTimestamps(t *testing.T) {
	var d Data
	err := json.Unmarshal([]byte(`{
		"id": "0GAJ7GWF5QBBM4S7YYHHCCC4F1",
		"value": "2",
		"created_at": "2026-11-24T16:05:37Z",
		"created_epoch": 1795536337,
		"expiration": "2027-01-23T16:05:37Z",
		"updated_at": null
	}`), &d)

	assert.NoError(t, err)
	assert.True(t, d.GetCreatedAt().Equal(Timestamp{time.Date(2026, 11, 24, 16, 5, 37, 0, time.UTC)}))
	assert.True(t, d.GetExpiration().Equal(Timestamp{time.Date(2027, 1, 23, 16, 5, 37, 0, time.UTC)}))
	assert.Nil(t, d.UpdatedAt)
	assert.True(t, d.GetUpdatedAt().IsZero())
}

// Single-record data endpoints (get, first, last, previous) return expiration
// as a quoted Unix timestamp, unlike the list and create endpoints.
func TestDataDecodesSingleRecordTimestamps(t *testing.T) {
	var d Data
	err := json.Unmarshal([]byte(`{
		"id": "0GAJ91ZR811KT6WZHRBH1G7VG0",
		"feed_id": 25,
		"value": "1",
		"location": null,
		"created_at": "2026-11-24T17:31:26Z",
		"updated_at": "2026-11-24T17:31:26Z",
		"expiration": "1800725486.0",
		"lat": null,
		"lon": null,
		"ele": null
	}`), &d)

	assert.NoError(t, err)
	assert.True(t, d.GetCreatedAt().Equal(Timestamp{time.Date(2026, 11, 24, 17, 31, 26, 0, time.UTC)}))
	assert.True(t, d.GetExpiration().Equal(Timestamp{time.Unix(1800725486, 0)}))
	assert.Equal(t, "2027-01-23T17:31:26Z", d.GetExpiration().Format(time.RFC3339))
}
