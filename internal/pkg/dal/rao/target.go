package rao

type CreateTargetReq struct {
	TargetID   string    `json:"target_id"`
	ParentID   string    `json:"parent_id"`
	ProjectID  string    `json:"project_id"`
	Mark       string    `json:"mark"`
	TargetType string    `json:"target_type"`
	Name       string    `json:"name"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Sort       int       `json:"sort"`
	TypeSort   string    `json:"type_sort"`
	Request    *Request  `json:"request"`
	Response   *Response `json:"response"`
	Mock       string    `json:"mock"`
	MockURL    string    `json:"mock_url"`
	Version    int       `json:"version"`
	IsChanged  int       `json:"is_changed"`
	IsSave     int       `json:"is_save"`
	IsForce    int       `json:"is_force"`
}

type CreateTargetResp struct {
}

type ListTargetReq struct {
}

type ListTargetResp struct {
}
