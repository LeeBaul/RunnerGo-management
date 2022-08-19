package rao

type AuthSignupReq struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
	Nickname       string `json:"nickname"`
	//VerifyCode     string `json:"verify_code"`
}

type AuthSignupResp struct {
	Token         string `json:"token"`
	ExpireTimeSec int64  `json:"expire_time_sec"`
}

type AuthLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResp struct {
	Token         string `json:"token"`
	ExpireTimeSec int64  `json:"expire_time_sec"`
}

type AuthSendMailVerifyReq struct {
	Email string `json:"email"`
}

type AuthSendMailVerifyResp struct {
}

type AuthUpdatePasswordReq struct {
}

type AuthUpdatePasswordResp struct {
}
