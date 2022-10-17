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

type ForgetPasswordReq struct {
	Email string `json:"email"`
}

type ForgetPasswordResp struct {
}

type AuthResetPasswordReq struct {
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type AuthResetPasswordResp struct {
}

type UpdatePasswordReq struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6,eqfield=RepeatPassword"`
	RepeatPassword  string `json:"repeat_password" binding:"required,min=6"`
}

type UpdatePasswordResp struct {
}

type UpdateNicknameReq struct {
	Nickname string `json:"nickname" binding:"required,min=2"`
}

type UpdateNicknameResp struct {
}

type UpdateAvatarReq struct {
	AvatarURL string `json:"avatar_url" binding:"required"`
}

type UpdateAvatarResp struct {
}

type VerifyPasswordReq struct {
	Password string `json:"password"`
}

type VerifyPasswordResp struct {
	IsMatch bool `json:"is_match"`
}

type ResetPasswordReq struct {
	U              int64  `json:"u"`
	NewPassword    string `json:"new_password" binding:"required,min=6,eqfield=RepeatPassword"`
	RepeatPassword string `json:"repeat_password" binding:"required,min=6"`
}

type ResetPasswordResp struct {
}
