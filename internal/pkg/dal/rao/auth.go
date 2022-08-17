package rao

type AuthSignupReq struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
	Nickname       string `json:"nickname"`
	//VerifyCode     string `json:"verify_code"`
}

type AuthSignupResp struct {
	User          *AuthUser `json:"user"`
	Token         string    `json:"token"`
	ExpireTimeSec int64     `json:"expire_time_sec"`
}

type AuthLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResp struct {
	User          *AuthUser `json:"user"`
	Token         string    `json:"token"`
	ExpireTimeSec int64     `json:"expire_time_sec"`
}

type AuthUser struct {
	Email    string      `json:"email"`
	Nickname string      `json:"nickname"`
	Teams    []*AuthTeam `json:"teams"`
}

type AuthTeam struct {
	TeamID   int64  `json:"team_id"`
	TeamName string `json:"team_name"`
	Sort     int32  `json:"sort"`
}
