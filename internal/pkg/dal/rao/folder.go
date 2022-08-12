package rao

type SaveFolderReq struct {
	TargetID   string   `json:"target_id"`
	ParentID   string   `json:"parent_id"`
	TargetType string   `json:"target_type"`
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Sort       int      `json:"sort"`
	TypeSort   int      `json:"type_sort"`
	Version    int      `json:"version"`
	Request    *Request `json:"request"`
	Script     *Script  `json:"script"`
}

type Request struct {
	Description string        `json:"description"`
	Header      *Header       `json:"header"`
	Query       []interface{} `json:"query"`
	Body        *Body         `json:"body"`
	Auth        *Auth         `json:"auth"`
}

type Script struct {
	PreScript       string `json:"pre_script"`
	Test            string `json:"test"`
	PreScriptSwitch bool   `json:"pre_script_switch"`
	TestSwitch      bool   `json:"test_switch"`
}
