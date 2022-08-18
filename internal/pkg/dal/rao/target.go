package rao

type CreateTargetReq struct {
	TargetID    int64     `json:"target_id"`
	ParentID    int64     `json:"parent_id"`
	ProjectID   string    `json:"project_id"`
	Mark        string    `json:"mark"`
	TargetType  string    `json:"target_type"`
	Name        string    `json:"name"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	Sort        int32     `json:"sort"`
	TypeSort    int32     `json:"type_sort"`
	Request     *Request  `json:"request"`
	Response    *Response `json:"response"`
	Version     int32     `json:"version"`
	Description string    `json:"description"`
}

type CreateTargetResp struct {
}

type DeleteTargetReq struct {
	TargetID int64 `json:"target_id"`
}

type DeleteTargetResp struct {
}

type ListTargetReq struct {
}

type ListTargetResp struct {
}
