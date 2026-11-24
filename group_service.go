// GroupService provides CRUD access to Groups.

package adafruitio

import "context"

// Group is an Adafruit IO group of feeds. See Feed for how pointer fields
// work.
type Group struct {
	ID          *int64     `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Key         *string    `json:"key,omitempty"`
	Owner       *Owner     `json:"owner,omitempty"`
	UserID      *int64     `json:"user_id,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
	Feeds       []*Feed    `json:"feeds,omitempty"`
	Visibility  *string    `json:"visibility,omitempty"`
	Shared      *bool      `json:"is_shared,omitempty"`
}

type GroupService struct {
	client *Client
}

// All returns all Groups for the current account.
func (s *GroupService) All(ctx context.Context) ([]*Group, *Response, error) {
	path := "groups"

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	groups := make([]*Group, 0)
	resp, err := s.client.Do(ctx, req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// Create makes a new Group and either returns a new Group instance or an error.
func (s *GroupService) Create(ctx context.Context, g *Group) (*Group, *Response, error) {
	path := "groups"

	req, rerr := s.client.NewRequest("POST", path, g)
	if rerr != nil {
		return nil, nil, rerr
	}

	var group Group
	resp, err := s.client.Do(ctx, req, &group)
	if err != nil {
		return nil, resp, err
	}

	return &group, resp, nil
}

// Get returns the Group identified by key.
func (s *GroupService) Get(ctx context.Context, key string) (*Group, *Response, error) {
	path, perr := groupPath(key)
	if perr != nil {
		return nil, nil, perr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var group Group
	resp, err := s.client.Do(ctx, req, &group)
	if err != nil {
		return nil, resp, err
	}

	return &group, resp, nil
}

// Update takes an ID and a Group record, updates it, and returns a new Group
// instance or an error.
func (s *GroupService) Update(ctx context.Context, key string, group *Group) (*Group, *Response, error) {
	path, perr := groupPath(key)
	if perr != nil {
		return nil, nil, perr
	}

	req, rerr := s.client.NewRequest("PATCH", path, group)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedGroup Group
	resp, err := s.client.Do(ctx, req, &updatedGroup)
	if err != nil {
		return nil, resp, err
	}

	return &updatedGroup, resp, nil
}

// Delete the Group identified by the given ID.
func (s *GroupService) Delete(ctx context.Context, key string) (*Response, error) {
	path, perr := groupPath(key)
	if perr != nil {
		return nil, perr
	}

	req, rerr := s.client.NewRequest("DELETE", path, nil)
	if rerr != nil {
		return nil, rerr
	}

	return s.client.Do(ctx, req, nil)
}

// Feeds returns the Feeds in the Group identified by key. Within a group
// other than the default group, feed keys are scoped to the group, for
// example "mygroup.temperature".
func (s *GroupService) Feeds(ctx context.Context, key string) ([]*Feed, *Response, error) {
	path, perr := groupPath(key, "feeds")
	if perr != nil {
		return nil, nil, perr
	}

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	feeds := make([]*Feed, 0)
	resp, err := s.client.Do(ctx, req, &feeds)
	if err != nil {
		return nil, resp, err
	}

	return feeds, resp, nil
}

// CreateFeed creates a Feed in the Group identified by key, and returns the
// created Feed.
func (s *GroupService) CreateFeed(ctx context.Context, key string, feed *Feed) (*Feed, *Response, error) {
	path, perr := groupPath(key, "feeds")
	if perr != nil {
		return nil, nil, perr
	}

	req, rerr := s.client.NewRequest("POST", path, feed)
	if rerr != nil {
		return nil, nil, rerr
	}

	var created Feed
	resp, err := s.client.Do(ctx, req, &created)
	if err != nil {
		return nil, resp, err
	}

	return &created, resp, nil
}

// AddFeed adds the Feed identified by feedKey to the Group identified by key,
// and returns the updated Group. The feed stays in its other groups.
func (s *GroupService) AddFeed(ctx context.Context, key, feedKey string) (*Group, *Response, error) {
	return s.changeFeed(ctx, key, "add", feedKey)
}

// RemoveFeed removes the Feed identified by feedKey from the Group identified
// by key, and returns the updated Group. The feed itself is not deleted.
func (s *GroupService) RemoveFeed(ctx context.Context, key, feedKey string) (*Group, *Response, error) {
	return s.changeFeed(ctx, key, "remove", feedKey)
}

// changeFeed handles the add and remove actions for AddFeed and RemoveFeed.
func (s *GroupService) changeFeed(ctx context.Context, key, action, feedKey string) (*Group, *Response, error) {
	path, perr := groupPath(key, action)
	if perr != nil {
		return nil, nil, perr
	}
	if feedKey == "" {
		return nil, nil, errEmptyFeedKey
	}
	body := struct {
		FeedKey string `json:"feed_key"`
	}{feedKey}

	req, rerr := s.client.NewRequest("POST", path, body)
	if rerr != nil {
		return nil, nil, rerr
	}

	var group Group
	resp, err := s.client.Do(ctx, req, &group)
	if err != nil {
		return nil, resp, err
	}

	return &group, resp, nil
}

// LastValues returns the most recent value of every Feed in the Group
// identified by key, keyed by feed key. Within a group other than the default
// group, feed keys are scoped to the group, for example
// "mygroup.temperature". A feed with no data has a nil value.
func (s *GroupService) LastValues(ctx context.Context, key string) (map[string]*string, *Response, error) {
	path, perr := groupPath(key, "data", "last_value")
	if perr != nil {
		return nil, nil, perr
	}
	path += "?type=json"

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var values struct {
		Feeds map[string]*string `json:"feeds"`
	}
	resp, err := s.client.Do(ctx, req, &values)
	if err != nil {
		return nil, resp, err
	}

	if values.Feeds == nil {
		values.Feeds = map[string]*string{}
	}
	return values.Feeds, resp, nil
}
