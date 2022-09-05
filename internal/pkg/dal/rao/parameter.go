package rao

type Auth struct {
	Type   string  `json:"type"`
	Kv     *KV     `json:"kv"`
	Bearer *Bearer `json:"bearer"`
	Basic  *Basic  `json:"basic"`
}

type Bearer struct {
	Key string `json:"key"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Basic struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Query struct {
	Parameter []*Parameter `json:"parameter"`
}

type Header struct {
	Parameter []*Parameter `json:"parameter"`
}

type Body struct {
	Mode      string       `json:"mode"`
	Parameter []*Parameter `json:"parameter"`
	Raw       string       `json:"raw"`
	RawPara   []*Parameter `json:"raw_para"`
}

type Parameter struct {
	IsChecked   int32  `json:"is_checked"`
	Type        string `json:"type"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	NotNull     int32  `json:"not_null"`
	Description string `json:"description"`
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
	Header      *Header `json:"header"`
	Query       *Query  `json:"query"`
	Event       *Event  `json:"event"`
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
