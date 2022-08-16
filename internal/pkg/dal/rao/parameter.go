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
	IsChecked   string `json:"is_checked"`
	Type        string `json:"type"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	NotNull     string `json:"not_null"`
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
