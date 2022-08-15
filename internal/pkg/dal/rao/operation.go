package rao

type ListOperationReq struct {
}

type ListOperationResp struct {
	Operations []*Operation `json:"operations"`
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
