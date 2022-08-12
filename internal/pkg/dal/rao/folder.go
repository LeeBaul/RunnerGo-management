package rao

type SaveFolderReq struct {
	TargetID   int64    `json:"target_id"`
	ParentID   int64    `json:"parent_id"`
	TargetType string   `json:"target_type"`
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Sort       int32    `json:"sort"`
	TypeSort   int32    `json:"type_sort"`
	Version    int32    `json:"version"`
	Request    *Request `json:"request"`
	Script     *Script  `json:"script"`
}

type Request struct {
	Description string       `json:"description"`
	Header      []*Parameter `json:"header"`
	Query       []*Parameter `json:"query"`
	Body        []*Parameter `json:"body"`
	Auth        *Auth        `json:"auth"`
}

type Script struct {
	PreScript       string `json:"pre_script"`
	Test            string `json:"test"`
	PreScriptSwitch bool   `json:"pre_script_switch"`
	TestSwitch      bool   `json:"test_switch"`
}
