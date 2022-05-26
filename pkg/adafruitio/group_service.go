// GroupService provides CRUD access to Groups.

package adafruitio

import "fmt"

type Group struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
	Source      string   `json:"source,omitempty"`
	SourceKeys  []string `json:"source_keys,omitempty"`
	Feeds       []*Feed  `json:"feeds,omitempty"`
	Visibility  string   `json:"visibility"`
}

type GroupService struct {
	client *Client
}

// All returns all Groups for the current account.
func (s *GroupService) All() ([]*Group, *Response, error) {
	path := "api/v1/groups"

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	// request populates Feed slice
	groups := make([]*Group, 0)
	resp, err := s.client.Do(req, &groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}

// Create makes a new Group and either returns a new Group instance or an error.
func (s *GroupService) Create(g *Group) (*Group, *Response, error) {
	path := "api/v1/groups"

	req, rerr := s.client.NewRequest("POST", path, g)
	if rerr != nil {
		return nil, nil, rerr
	}

	var group Group
	resp, err := s.client.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}

	return &group, resp, nil
}

// Get returns the Group record identified by the given ID
func (s *GroupService) Get(id interface{}) (*Group, *Response, error) {
	path := fmt.Sprintf("api/v1/groups/%v", id)

	req, rerr := s.client.NewRequest("GET", path, nil)
	if rerr != nil {
		return nil, nil, rerr
	}

	var group Group
	resp, err := s.client.Do(req, &group)
	if err != nil {
		return nil, resp, err
	}

	return &group, resp, nil
}

// Update takes an ID and a Group record, updates it, and returns a new Group
// instance or an error.
func (s *GroupService) Update(id interface{}, group *Group) (*Group, *Response, error) {
	path := fmt.Sprintf("api/v1/groups/%v", id)

	req, rerr := s.client.NewRequest("PATCH", path, group)
	if rerr != nil {
		return nil, nil, rerr
	}

	var updatedGroup Group
	resp, err := s.client.Do(req, &updatedGroup)
	if err != nil {
		return nil, resp, err
	}

	return &updatedGroup, resp, nil
}

// Delete the Group identified by the given ID.
func (s *GroupService) Delete(id interface{}) (*Response, error) {
	path := fmt.Sprintf("api/v1/groups/%v", id)

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
