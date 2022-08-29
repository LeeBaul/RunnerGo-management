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
	SceneID int64   `json:"scene_id"`
	TeamID  int64   `json:"team_id"`
	Version int32   `json:"version"`
	Flows   []*Flow `json:"flows"`
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

// 接口详情
type API struct {
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
