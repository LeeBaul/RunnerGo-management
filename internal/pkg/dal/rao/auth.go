package rao

type AuthSignupReq struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeat_password" binding:"required,eqfiel=password"`
	Nickname       string `json:"nickname" binding:"required"`
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

type SetUserSettingsReq struct {
	UserID       int64        `json:"user_id"`
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
	CurrentTeamID int64 `json:"current_team_id"`
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
