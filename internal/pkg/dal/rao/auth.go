package rao

type AuthSignupReq struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
	Nickname       string `json:"nickname"`
	VerifyCode     string `json:"verify_code"`
}

type AuthSignupResp struct {
}

type AuthLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResp struct {
}
