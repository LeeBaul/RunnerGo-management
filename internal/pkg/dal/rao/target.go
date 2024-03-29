package rao

type SendTargetReq struct {
	TargetID int64 `json:"target_id" binding:"required,gt=0"`
	TeamID   int64 `json:"team_id" binding:"required,gt=0"`
}

type SendTargetResp struct {
	RetID string `json:"ret_id"`
}

type GetSendTargetResultReq struct {
	RetID string `form:"ret_id" binding:"required,gt=0"`
}

type GetSendTargetResultResp struct {
}

type SaveTargetReq struct {
	TargetID    int64     `json:"target_id"`
	ParentID    int64     `json:"parent_id"`
	TeamID      int64     `json:"team_id" binding:"required,gt=0"`
	Mark        string    `json:"mark"`
	Name        string    `json:"name" binding:"required,min=1"`
	Method      string    `json:"method" binding:"required"`
	URL         string    `json:"url"`
	Sort        int32     `json:"sort"`
	TypeSort    int32     `json:"type_sort"`
	Request     *Request  `json:"request"`
	Response    *Response `json:"response"`
	Version     int32     `json:"version"`
	Description string    `json:"description"`
	Assert      []*Assert `json:"assert"`
	Regex       []*Regex  `json:"regex"`
}

type SaveTargetResp struct {
	TargetID int64 `json:"target_id"`
}

type SortTargetReq struct {
	Targets []*SortTarget `json:"targets"`
}

type SortTarget struct {
	TeamID   int64 `json:"team_id"`
	TargetID int64 `json:"target_id"`
	Sort     int32 `json:"sort"`
	ParentID int64 `json:"parent_id"`
}

type SortTargetResp struct {
}

type TrashTargetReq struct {
	TargetID int64 `json:"target_id" binding:"required,gt=0"`
}

type TrashTargetResp struct {
}

type RecallTargetReq struct {
	TargetID int64 `json:"target_id" binding:"required,gt=0"`
}

type RecallTargetResp struct {
}

type DeleteTargetReq struct {
	TargetID int64 `json:"target_id" binding:"required,gt=0"`
}

type DeleteTargetResp struct {
}

type ListTrashTargetReq struct {
	TeamID int64 `form:"team_id" binding:"required,gt=0"`
	Page   int   `form:"page,default=1"`
	Size   int   `form:"size,default=10"`
}

type ListTrashTargetResp struct {
	Targets []*FolderAPI `json:"targets"`
	Total   int64        `json:"total"`
}

type ListFolderAPIReq struct {
	TeamID int64 `form:"team_id" binding:"required,gt=0"`
	Page   int   `form:"page,default=1"`
	Size   int   `form:"size,default=10"`
}

type ListFolderAPIResp struct {
	Targets []*FolderAPI `json:"targets"`
	Total   int64        `json:"total"`
}

type FolderAPI struct {
	TargetID      int64  `json:"target_id"`
	TeamID        int64  `json:"team_id"`
	TargetType    string `json:"target_type"`
	Name          string `json:"name"`
	ParentID      int64  `json:"parent_id"`
	Method        string `json:"method"`
	Sort          int32  `json:"sort"`
	TypeSort      int32  `json:"type_sort"`
	Version       int32  `json:"version"`
	CreatedUserID int64  `json:"created_user_id"`
	RecentUserID  int64  `json:"recent_user_id"`
}

type ListGroupSceneReq struct {
	TeamID int64 `form:"team_id" binding:"required,gt=0"`
	Source int32 `form:"source,default=1"`
	PlanID int64 `form:"plan_id"`
	Page   int   `form:"page,default=1"`
	Size   int   `form:"size,default=10"`
}

type ListGroupSceneResp struct {
	Targets []*GroupScene `json:"targets"`
	Total   int64         `json:"total"`
}

type GroupScene struct {
	TargetID      int64  `json:"target_id"`
	TeamID        int64  `json:"team_id"`
	TargetType    string `json:"target_type"`
	Name          string `json:"name"`
	ParentID      int64  `json:"parent_id"`
	Method        string `json:"method"`
	Sort          int32  `json:"sort"`
	TypeSort      int32  `json:"type_sort"`
	Version       int32  `json:"version"`
	CreatedUserID int64  `json:"created_user_id"`
	RecentUserID  int64  `json:"recent_user_id"`
	Description   string `json:"description"`
}

type BatchGetDetailReq struct {
	TeamID    int64   `form:"team_id" binding:"required,gt=0"`
	TargetIDs []int64 `form:"target_ids" binding:"required,gt=0"`
}

type BatchGetDetailResp struct {
	Targets []*APIDetail `json:"targets"`
}

type APIDetail struct {
	TargetID       int64         `json:"target_id"`
	ParentID       int64         `json:"parent_id"`
	TargetType     string        `json:"target_type"`
	TeamID         int64         `json:"team_id"`
	Name           string        `json:"name"`
	Method         string        `json:"method"`
	URL            string        `json:"url"`
	Sort           int32         `json:"sort"`
	TypeSort       int32         `json:"type_sort"`
	Request        *Request      `json:"request"`
	Response       *Response     `json:"response"`
	Version        int32         `json:"version"`
	Description    string        `json:"description"`
	CreatedTimeSec int64         `json:"created_time_sec"`
	UpdatedTimeSec int64         `json:"updated_time_sec"`
	Assert         []*Assert     `json:"assert"`
	Regex          []*Regex      `json:"regex"`
	Variable       []*KVVariable `json:"variable"` // 全局变量
}

type KVVariable struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
