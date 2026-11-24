// Adapted from the Timestamp type in the go-github project located at
// https://github.com/google/go-github

package adafruitio

import (
	"math"
	"strconv"
	"time"
)

// Timestamp represents a time returned by the Adafruit IO API, such as a
// record's created_at. It embeds time.Time, so all of time.Time's methods are
// available.
//
// Timestamp decodes RFC 3339 strings, such as "2026-11-24T16:05:19Z", as well
// as Unix timestamps in seconds, quoted or not. It encodes as an RFC 3339
// string.
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// GetTime returns a pointer to the underlying time.Time, or nil if t is nil.
func (t *Timestamp) GetTime() *time.Time {
	if t == nil {
		return nil
	}
	return &t.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface. Times are expected
// as RFC 3339 strings, or as Unix timestamps in seconds, which may be quoted
// and may have a fractional part. Some API endpoints return a record's
// expiration as a quoted timestamp such as "1800725486.0".
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}

	if len(str) > 0 && str[0] == '"' {
		s, err := strconv.Unquote(str)
		if err != nil {
			return err
		}
		if s == "" {
			t.Time = time.Time{}
			return nil
		}
		parsed, err := time.Parse(time.RFC3339, s)
		if err != nil {
			// fall back to a quoted Unix timestamp
			if unix, uerr := parseUnixSeconds(s); uerr == nil {
				t.Time = unix
				return nil
			}
			return err
		}
		t.Time = parsed
		return nil
	}

	unix, err := parseUnixSeconds(str)
	if err != nil {
		return err
	}
	t.Time = unix
	return nil
}

// parseUnixSeconds parses a Unix timestamp in seconds, such as "1800725486"
// or "1800725486.5", as a UTC time.
func parseUnixSeconds(s string) (time.Time, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return time.Time{}, err
	}
	sec, frac := math.Modf(f)
	return time.Unix(int64(sec), int64(math.Round(frac*1e9))).UTC(), nil
}

// Equal reports whether t and u are equal based on time.Time.Equal.
func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}
