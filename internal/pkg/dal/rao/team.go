package rao

type ListTeamReq struct {
}

type ListTeamResp struct {
	Teams []*Team `json:"teams"`
}

type Team struct {
	Name   string `json:"name"`
	Sort   int32  `json:"sort"`
	TeamID int64  `json:"team_id"`
	RoleID int64  `json:"role_id"`
}

type ListMembersReq struct {
	TeamID int64 `form:"team_id"`
}

type ListMembersResp struct {
	Members []*Member `json:"members"`
}

type Member struct {
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname"`
	JoinTimeSec int64  `json:"join_time_sec"`
}
