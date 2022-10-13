package rao

type AuthSignupReq struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6,eqfield=RepeatPassword"`
	RepeatPassword string `json:"repeat_password" binding:"required,min=6"`
	Nickname       string `json:"nickname" binding:"required,min=2"`
}

type AuthSignupResp struct {
	Token         string `json:"token"`
	ExpireTimeSec int64  `json:"expire_time_sec"`
}

type AuthLoginReq struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6,max=32"`
	IsAutoLogin bool   `json:"is_auto_login"`
}

type AuthLoginResp struct {
	Token         string `json:"token"`
	ExpireTimeSec int64  `json:"expire_time_sec"`
	IsAPIPostUser bool   `json:"is_api_post_user"`
}

type SetUserSettingsReq struct {
	UserSettings UserSettings `json:"settings"`
}

type SetUserSettingsResp struct {
}

type GetUserSettingsReq struct {
}

type GetUserSettingsResp struct {
	UserSettings *UserSettings `json:"settings"`
}

type UserSettings struct {
	CurrentTeamID int64 `json:"current_team_id" binding:"required,gt=0"`
}

type AuthSendMailVerifyReq struct {
	Email string `json:"email"`
}

type AuthSendMailVerifyResp struct {
}

type AuthResetPasswordReq struct {
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type AuthResetPasswordResp struct {
}
