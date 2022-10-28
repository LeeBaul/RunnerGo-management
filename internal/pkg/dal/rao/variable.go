package rao

type SaveVariableReq struct {
	VarID       int64  `json:"var_id"`
	TeamID      int64  `json:"team_id" binding:"required,gt=0"`
	Var         string `json:"var" binding:"required,min=1"`
	Val         string `json:"val" binding:"required,min=1"`
	Description string `json:"description"`
}

type SaveVariableResp struct {
}

type DeleteVariableReq struct {
	TeamID int64 `json:"team_id"`
	VarID  int64 `json:"var_id"`
}

type DeleteVariableResp struct {
}

type SyncVariablesReq struct {
	TeamID    int64       `json:"team_id" binding:"required,gt=0"`
	Variables []*Variable `json:"variables"`
}

type SyncVariablesResp struct {
}

type SyncSceneVariablesReq struct {
	TeamID    int64       `json:"team_id" binding:"required,gt=0"`
	SceneID   int64       `json:"scene_id" binding:"required,gt=0"`
	Variables []*Variable `json:"variables"`
}

type SyncSceneVariablesResp struct {
}

type ListSceneVariablesReq struct {
	TeamID  int64 `form:"team_id" binding:"required,gt=0"`
	SceneID int64 `form:"scene_id" binding:"required,gt=0"`
	Page    int   `form:"page"`
	Size    int   `form:"size"`
}

type ListVariablesReq struct {
	TeamID int64 `form:"team_id" binding:"required,gt=0"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

type ListVariablesResp struct {
	Variables []*Variable `json:"variables"`
	Total     int64       `json:"total"`
}

type Variable struct {
	VarID       int64  `json:"var_id,omitempty"`
	TeamID      int64  `json:"team_id,omitempty"`
	Var         string `json:"var"`
	Val         string `json:"val"`
	Description string `json:"description"`
}

type ImportVariablesReq struct {
	TeamID  int64  `json:"team_id" binding:"required,gt=0"`
	SceneID int64  `json:"scene_id" binding:"required,gt=0"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

type ImportVariablesResp struct {
}

type DeleteImportSceneVariablesReq struct {
	TeamID  int64  `json:"team_id" binding:"required,gt=0"`
	SceneID int64  `json:"scene_id" binding:"required,gt=0"`
	Name    string `json:"name"`
}

type ListImportVariablesReq struct {
	TeamID  int64 `form:"team_id" binding:"required,gt=0"`
	SceneID int64 `form:"scene_id" binding:"required,gt=0"`
}

type ListImportVariablesResp struct {
	Imports []*Import `json:"imports"`
}

type Import struct {
	TeamID         int64  `json:"team_id"`
	SceneID        int64  `json:"scene_id"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	CreatedTimeSec int64  `json:"created_time_sec"`
}
