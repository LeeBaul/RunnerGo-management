package rao

type SaveGroupReq struct {
	TeamID      int64  `json:"team_id" binding:"required,gt=0"`
	TargetID    int64  `json:"target_id"`
	ParentID    int64  `json:"parent_id"`
	Name        string `json:"name" binding:"required,max=32"`
	Method      string `json:"method"`
	Sort        int32  `json:"sort"`
	TypeSort    int32  `json:"type_sort"`
	Version     int32  `json:"version" binding:"required,gt=0"`
	Source      int32  `json:"source"`
	PlanID      int64  `json:"plan_id"`
	Description string `json:"description"`
	//Request  *Request `json:"request"`
	//Script   *Script  `json:"script"`
}

type SaveGroupResp struct {
}

type GetGroupReq struct {
	TeamID   int64 `form:"team_id"`
	TargetID int64 `form:"target_id"`
}

type GetGroupResp struct {
	Group *Group `json:"group"`
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
	Source   int32    `json:"source"`
	PlanID   int64    `json:"plan_id"`
}
