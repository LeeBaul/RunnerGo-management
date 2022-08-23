package rao

type ListOperationReq struct {
	TeamID int64 `form:"team_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

type ListOperationResp struct {
	Operations []*Operation `json:"operations"`
	Total      int64        `json:"total"`
}

type Operation struct {
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	UserAvatar     string `json:"user_avatar"`
	UserStatus     int32  `json:"user_status"`
	Category       int32  `json:"category"`
	Name           string `json:"name"`
	CreatedTimeSec int64  `json:"created_time_sec"`
}
