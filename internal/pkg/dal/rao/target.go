package rao

type CreateTargetReq struct {
	Auth        *Auth   `json:"auth"`
	Header      *Header `json:"header"`
	Body        *Body   `json:"body"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
}

type CreateTargetResp struct {
}

type ListTargetReq struct {
}

type ListTargetResp struct {
}
