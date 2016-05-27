package adafruitio

/*
[
  {
    "id": 0,
    "name": "string",
    "description": "string",
    "feeds": [
      {
        "id": 0,
        "name": "string",
        "key": "string",
        "description": "string",
        "unit_type": "string",
        "unit_symbol": "string",
        "history": true,
        "visibility": "private",
        "license": "string",
        "enabled": true,
        "last_value": "string",
        "status": "string",
        "group_id": 0,
        "created_at": "string",
        "updated_at": "string"
      }
    ],
    "created_at": "string",
    "updated_at": "string"
  }
]
*/

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

// Get all Groups for the current account.
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
		return nil, nil, err
	}

	return groups, resp, nil
}

// Create a new Group
func (s *GroupService) Create(g *Group) (*Group, *Response, error) {
	path := "api/v1/groups"

	req, rerr := s.client.NewRequest("POST", path, g)
	if rerr != nil {
		return nil, nil, rerr
	}

	var group Group
	resp, err := s.client.Do(req, &group)
	if err != nil {
		return nil, nil, err
	}

	return &group, resp, nil
}
