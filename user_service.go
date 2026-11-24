package adafruitio

import (
	"context"
	"encoding/json"
)

// User is an Adafruit IO account. See Feed for how pointer fields work.
type User struct {
	ID                 *int64     `json:"id,omitempty"`
	Name               *string    `json:"name,omitempty"`
	Username           *string    `json:"username,omitempty"`
	Color              *string    `json:"color,omitempty"`
	Role               *string    `json:"role,omitempty"`
	TimeZone           *string    `json:"time_zone,omitempty"`
	DefaultGroupID     *int64     `json:"default_group_id,omitempty"`
	DefaultDashboardID *int64     `json:"default_dashboard_id,omitempty"`
	CreatedAt          *Timestamp `json:"created_at,omitempty"`
	UpdatedAt          *Timestamp `json:"updated_at,omitempty"`
}

// Throttle reports an account's rate limits and current usage. Rates are
// per minute.
type Throttle struct {
	// DataRateLimit is the number of data points the account may send per
	// minute, and ActiveDataRate is the number sent in the last minute.
	DataRateLimit  *int `json:"data_rate_limit,omitempty"`
	ActiveDataRate *int `json:"active_data_rate,omitempty"`

	ActiveSMSRate              *int    `json:"active_sms_rate,omitempty"`
	AuthenticationRate         *int    `json:"authentication_rate,omitempty"`
	SubscribeAuthorizationRate *int    `json:"subscribe_authorization_rate,omitempty"`
	PublishAuthorizationRate   *int    `json:"publish_authorization_rate,omitempty"`
	HourlyBanRate              *int    `json:"hourly_ban_rate,omitempty"`
	MQTTBanErrorMessage        *string `json:"mqtt_ban_error_message,omitempty"`
}

// UserService provides access to the current user's account.
type UserService struct {
	client *Client
}

// Get returns the user that owns the client's API key.
func (s *UserService) Get(ctx context.Context) (*User, *Response, error) {
	// /api/v2/user isn't under the username, so step up out of it
	req, rerr := s.client.NewRequest("GET", "../user", nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var raw json.RawMessage
	resp, err := s.client.Do(ctx, req, &raw)
	if err != nil {
		return nil, resp, err
	}

	user, err := decodeUser(raw)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// decodeUser decodes a /user response. The API documents a plain user record,
// but returns one wrapped in a "user" field alongside other account details,
// so both forms are accepted.
func decodeUser(raw json.RawMessage) (*User, error) {
	var wrapped struct {
		User *User `json:"user"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.User != nil {
		return wrapped.User, nil
	}

	var user User
	if err := json.Unmarshal(raw, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Throttle returns the account's rate limits and current usage.
func (s *UserService) Throttle(ctx context.Context) (*Throttle, *Response, error) {
	req, rerr := s.client.NewRequest("GET", "throttle", nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var throttle Throttle
	resp, err := s.client.Do(ctx, req, &throttle)
	if err != nil {
		return nil, resp, err
	}

	return &throttle, resp, nil
}
