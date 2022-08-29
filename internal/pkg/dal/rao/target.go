package rao

type CreateTargetReq struct {
	TargetID    int64     `json:"target_id"`
	ParentID    int64     `json:"parent_id"`
	TeamID      int64     `json:"team_id"`
	ProjectID   string    `json:"project_id"`
	Mark        string    `json:"mark"`
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

type TrashTargetReq struct {
	TargetID int64 `json:"target_id"`
}

type TrashTargetResp struct {
}

type RecallTargetReq struct {
	TargetID int64 `json:"target_id"`
}

type RecallTargetResp struct {
}

type DeleteTargetReq struct {
	TargetID int64 `json:"target_id"`
}

type DeleteTargetResp struct {
}

type ListTrashTargetReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

type ListTrashTargetResp struct {
	Targets []*FolderAPI `json:"targets"`
	Total   int64        `json:"total"`
}

type ListFolderAPIReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
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
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
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
	TeamID    int64   `form:"team_id"`
	TargetIDs []int64 `form:"target_ids"`
}

type BatchGetDetailResp struct {
	Targets []*APIDetail `json:"targets"`
}

type APIDetail struct {
	TargetID    int64     `json:"target_id"`
	ParentID    int64     `json:"parent_id"`
	TeamID      int64     `json:"team_id"`
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
