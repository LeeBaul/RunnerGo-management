package rao

type SaveFolderReq struct {
	TargetID int64  `json:"target_id"`
	TeamID   int64  `json:"team_id" binding:"required,gt=0"`
	ParentID int64  `json:"parent_id"`
	Name     string `json:"name" binding:"required,max=32"`
	Method   string `json:"method"`
	Sort     int32  `json:"sort"`
	TypeSort int32  `json:"type_sort"`
	Version  int32  `json:"version" binding:"required,gt=0"`
	//Request  *Request `json:"request"`
	//Script   *Script  `json:"script"`
}

type SaveFolderResp struct {
}

type GetFolderReq struct {
	TeamID   int64 `form:"team_id"`
	TargetID int64 `form:"target_id"`
}

type GetFolderResp struct {
	Folder *Folder `json:"folder"`
}

type Folder struct {
	TargetID int64    `json:"target_id"`
	TeamID   int64    `json:"team_id"`
	ParentID int64    `json:"parent_id"`
	Name     string   `json:"name"`
	Method   string   `json:"method"`
	Sort     int32    `json:"sort"`
	TypeSort int32    `json:"type_sort"`
	Version  int32    `json:"version"`
	Request  *Request `json:"request"`
	Script   *Script  `json:"script"`
}
