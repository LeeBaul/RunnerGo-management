package rao

type SaveSceneReq struct {
	TeamID   int64    `json:"team_id" binding:"required,gt=0"`
	TargetID int64    `json:"target_id"`
	ParentID int64    `json:"parent_id"`
	Name     string   `json:"name" binding:"required,min=4,max=32"`
	Method   string   `json:"method"`
	Sort     int32    `json:"sort"`
	TypeSort int32    `json:"type_sort"`
	Version  int32    `json:"version"`
	Request  *Request `json:"request"`
	Script   *Script  `json:"script"`
}

type SaveSceneResp struct {
}

type GetSceneReq struct {
	TeamID   int64   `form:"team_id" binding:"required,gt=0"`
	TargetID []int64 `form:"target_id" binding:"required,gt=0"`
}

type GetSceneResp struct {
	Scenes []*Scene `json:"scenes"`
}

type Scene struct {
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

type SaveFlowReq struct {
	SceneID int64 `json:"scene_id" binding:"required,gt=0"`
	TeamID  int64 `json:"team_id" binding:"required,gt=0"`
	Version int32 `json:"version"`

	Nodes           []*Node `json:"nodes"`
	Edges           []*Edge `json:"edges"`
	MultiLevelNodes []byte  `json:"multi_level_nodes"`
}

type Edge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type SaveFlowResp struct {
}

type Node struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	IsCheck bool   `json:"is_check"`

	// 接口
	Weight            int        `json:"weight,omitempty"`
	Mode              int        `json:"mode,omitempty"`
	ErrorThreshold    int        `json:"error_threshold,omitempty"`
	ResponseThreshold int        `json:"response_threshold,omitempty"`
	RequestThreshold  int        `json:"request_threshold,omitempty"`
	PercentAge        int        `json:"percent_age,omitempty"`
	API               *APIDetail `json:"api,omitempty"`

	// 全局断言
	Assets []string `json:"assets,omitempty"`

	// 等待控制器
	WaitMs int `json:"wait_ms,omitempty"`

	// 条件控制器
	Var     string `json:"var,omitempty"`
	Compare string `json:"compare,omitempty"`
	Val     string `json:"val,omitempty"`
}

// API 接口详情
//type API struct {
//	TargetID    int64     `json:"target_id"`
//	ParentID    int64     `json:"parent_id"`
//	TeamID      int64     `json:"team_id"`
//	ProjectID   string    `json:"project_id"`
//	Mark        string    `json:"mark"`
//	Name        string    `json:"name"`
//	Method      string    `json:"method"`
//	URL         string    `json:"url"`
//	Request     *Request  `json:"request"`
//	Response    *Response `json:"response,omitempty"`
//	Version     int32     `json:"version"`
//	Description string    `json:"description"`
//}

type GetFlowReq struct {
	SceneID int64 `form:"scene_id" binding:"required,gt=0"`
	TeamID  int64 `form:"team_id" binding:"required,gt=0"`
}

type GetFlowResp struct {
	SceneID int64 `json:"scene_id"`
	TeamID  int64 `json:"team_id"`
	Version int32 `json:"version"`

	Nodes           []*Node `json:"nodes"`
	Edges           []*Edge `json:"edges"`
	MultiLevelNodes []byte  `json:"multi_level_nodes"`
}
