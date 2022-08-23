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
	RoleID      int64  `json:"role_id"`
	JoinTimeSec int64  `json:"join_time_sec,omitempty"`
}

type InviteMemberReq struct {
	TeamID      int64  `json:"team_id"`
	MemberEmail string `json:"member_email" binding:"email"`
}

type InviteMemberResp struct {
}

type RemoveMemberReq struct {
	TeamID   int64 `json:"team_id"`
	MemberID int64 `json:"member_id"`
}

type RemoveMemberResp struct {
}
