package adafruitio

import (
	"encoding/json"
	"strconv"
)

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// The API returns a record's lat, lon, and ele as numbers from some endpoints
// (create and batch) and as strings, such as "40.7", from others (list, get,
// first, last, previous, and last_value). Both forms are accepted.
func (d *Data) UnmarshalJSON(b []byte) error {
	type data Data // has Data's fields, but not this method
	aux := struct {
		*data
		// these shadow the embedded fields with the same JSON names
		Latitude  flexibleFloat `json:"lat"`
		Longitude flexibleFloat `json:"lon"`
		Elevation flexibleFloat `json:"ele"`
	}{data: (*data)(d)}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// like encoding/json, leave fields that aren't in the JSON unchanged
	if aux.Latitude.present {
		d.Latitude = aux.Latitude.value
	}
	if aux.Longitude.present {
		d.Longitude = aux.Longitude.value
	}
	if aux.Elevation.present {
		d.Elevation = aux.Elevation.value
	}
	return nil
}

// flexibleFloat decodes a JSON number, a string containing a number, or null.
type flexibleFloat struct {
	present bool     // the field was in the JSON
	value   *float64 // nil for null or ""
}

func (f *flexibleFloat) UnmarshalJSON(b []byte) error {
	f.present = true

	s := string(b)
	if s == "null" {
		return nil
	}
	if len(s) > 0 && s[0] == '"' {
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			return err
		}
		if unquoted == "" {
			return nil
		}
		s = unquoted
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	f.value = &v
	return nil
}
