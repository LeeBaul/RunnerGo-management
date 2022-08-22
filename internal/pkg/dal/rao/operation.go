package rao

type ListOperationReq struct {
	TeamID int64 `query:"team_id"`
	Page   int   `query:"page"`
	Size   int   `query:"size"`
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
