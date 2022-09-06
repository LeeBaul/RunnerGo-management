package rao

type SaveVariableReq struct {
	VarID       int64  `json:"var_id"`
	TeamID      int64  `json:"team_id" binding:"required,gt=0"`
	Var         string `json:"var" binding:"required"`
	Val         string `json:"val" binding:"required"`
	Description string `json:"description"`
}

type SaveVariableResp struct {
}

type ListVariablesReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

type ListVariablesResp struct {
	Variables []*Variable `json:"variables"`
	Total     int64       `json:"total"`
}

type DeleteVariableReq struct {
	TeamID int64 `json:"team_id"`
	VarID  int64 `json:"var_id"`
}

type DeleteVariableResp struct {
}

type Variable struct {
	VarID       int64  `json:"var_id"`
	TeamID      int64  `json:"team_id"`
	Var         string `json:"var"`
	Val         string `json:"val"`
	Description string `json:"description"`
}
