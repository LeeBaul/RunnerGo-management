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
