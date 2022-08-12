package rao

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
