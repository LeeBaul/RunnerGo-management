package rao

type SaveSceneReq struct {
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

type SaveSceneResp struct {
}

type SaveFlowReq struct {
	TeamID   int64 `json:"team_id"`
	TargetID int64 `json:"target_id"`
}

type SaveFlowResp struct {
}
