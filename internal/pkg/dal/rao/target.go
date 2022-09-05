package rao

type SaveTargetReq struct {
	TargetID    int64     `json:"target_id"`
	ParentID    int64     `json:"parent_id"`
	TeamID      int64     `json:"team_id" binding:"required,gt=0"`
	Mark        string    `json:"mark"`
	Name        string    `json:"name" binding:"required,gt=0"`
	Method      string    `json:"method" binding:"required"`
	URL         string    `json:"url" binding:"required"`
	Sort        int32     `json:"sort"`
	TypeSort    int32     `json:"type_sort"`
	Request     *Request  `json:"request"`
	Response    *Response `json:"response"`
	Version     int32     `json:"version"`
	Description string    `json:"description"`
	Assert      []*Assert `json:"assert"`
	Regex       []*Regex  `json:"regex"`
}

type Assert struct {
	ResponseType int32  `json:"response_type"`
	Var          string `json:"var"`
	Compare      string `json:"compare"`
	Val          string `json:"val"`
}

type Regex struct {
	Var     string `json:"var"`
	Express string `json:"express"`
	Val     string `json:"val"`
}

type SaveTargetResp struct {
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
}

type BatchGetDetailReq struct {
	TeamID    int64   `form:"team_id" binding:"required,gt=0"`
	TargetIDs []int64 `form:"target_ids" binding:"required,gt=0"`
}

type BatchGetDetailResp struct {
	Targets []*APIDetail `json:"targets"`
}

type APIDetail struct {
	TargetID       int64     `json:"target_id"`
	ParentID       int64     `json:"parent_id"`
	TargetType     string    `json:"target_type"`
	TeamID         int64     `json:"team_id"`
	Name           string    `json:"name"`
	Method         string    `json:"method"`
	URL            string    `json:"url"`
	Sort           int32     `json:"sort"`
	TypeSort       int32     `json:"type_sort"`
	Request        *Request  `json:"request"`
	Response       *Response `json:"response"`
	Version        int32     `json:"version"`
	Description    string    `json:"description"`
	CreatedTimeSec int64     `json:"created_time_sec"`
	UpdatedTimeSec int64     `json:"updated_time_sec"`
	Assert         []*Assert `json:"assert"`
	Regex          []*Regex  `json:"regex"`
}
