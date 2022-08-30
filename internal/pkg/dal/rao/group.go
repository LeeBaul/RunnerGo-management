package rao

type SaveGroupReq struct {
	TeamID   int64    `json:"team_id"`
	TargetID int64    `json:"target_id"`
	ParentID int64    `json:"parent_id"`
	Name     string   `json:"name"`
	Method   string   `json:"method"`
	Sort     int32    `json:"sort"`
	TypeSort int32    `json:"type_sort"`
	Version  int32    `json:"version"`
	Request  *Request `json:"request"`
	Script   *Script  `json:"script"`
}

type SaveGroupResp struct {
}

type GetGroupReq struct {
	TeamID   int64 `form:"team_id"`
	TargetID int64 `form:"target_id"`
}

type GetGroupResp struct {
}

type Group struct {
	TeamID   int64    `json:"team_id"`
	TargetID int64    `json:"target_id"`
	ParentID int64    `json:"parent_id"`
	Name     string   `json:"name"`
	Method   string   `json:"method"`
	Sort     int32    `json:"sort"`
	TypeSort int32    `json:"type_sort"`
	Version  int32    `json:"version"`
	Request  *Request `json:"request"`
	Script   *Script  `json:"script"`
}
