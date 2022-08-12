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

type Script struct {
	PreScript       string `json:"pre_script"`
	Test            string `json:"test"`
	PreScriptSwitch bool   `json:"pre_script_switch"`
	TestSwitch      bool   `json:"test_switch"`
}

type Event struct {
	PreScript string `json:"pre_script"`
	Test      string `json:"test"`
}

type Query struct {
	Parameter []*Parameter `json:"parameter"`
}

type Cookie struct {
	Parameter []*Parameter `json:"parameter"`
}

type Resful struct {
	Parameter []*Parameter `json:"parameter"`
}

type Request struct {
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Auth        *Auth   `json:"auth"`
	Body        *Body   `json:"body"`
	Event       *Event  `json:"event"`
	Header      *Header `json:"header"`
	Query       *Query  `json:"query"`
	Cookie      *Cookie `json:"cookie"`
	Resful      *Resful `json:"resful"`
}

type Success struct {
	Raw       string       `json:"raw"`
	Parameter []*Parameter `json:"parameter"`
}

type Error struct {
	Raw       string       `json:"raw"`
	Parameter []*Parameter `json:"parameter"`
}

type Response struct {
	Success *Success `json:"success"`
	Error   *Error   `json:"error"`
}
