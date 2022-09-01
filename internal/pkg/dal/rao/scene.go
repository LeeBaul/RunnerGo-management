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
	SceneID int64   `json:"scene_id" binding:"required,gt=0"`
	TeamID  int64   `json:"team_id" binding:"required,gt=0"`
	Version int32   `json:"version"`
	Flows   []*Flow `json:"flows" bson:"flows"`
}

type SaveFlowResp struct {
}

type Flow struct {
	Index     int    `json:"index"`
	FlowType  string `json:"flow_type"`
	NextIndex int    `json:"next_index,omitempty"`

	// 接口
	Weight            int  `json:"weight,omitempty"`
	Mode              int  `json:"mode,omitempty"`
	ErrorThreshold    int  `json:"error_threshold,omitempty"`
	ResponseThreshold int  `json:"response_threshold,omitempty"`
	RequestThreshold  int  `json:"request_threshold,omitempty"`
	API               *API `json:"api,omitempty"`

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
type API struct {
	TargetID    int64     `json:"target_id"`
	ParentID    int64     `json:"parent_id"`
	TeamID      int64     `json:"team_id"`
	ProjectID   string    `json:"project_id"`
	Mark        string    `json:"mark"`
	Name        string    `json:"name"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	Request     *Request  `json:"request"`
	Response    *Response `json:"response,omitempty"`
	Version     int32     `json:"version"`
	Description string    `json:"description"`
}

type GetFlowReq struct {
	SceneID int64 `form:"scene_id" binding:"required,gt=0"`
	TeamID  int64 `form:"team_id" binding:"required,gt=0"`
}

type GetFlowResp struct {
	SceneID int64   `json:"scene_id"`
	TeamID  int64   `json:"team_id"`
	Version int32   `json:"version"`
	Flows   []*Flow `json:"flows"`
}
