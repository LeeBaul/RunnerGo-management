package team

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
)

func ListByUserID(ctx context.Context, userID int64) ([]*model.Team, error) {
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, err
	}

	var teamIDs []int64
	for _, team := range userTeams {
		teamIDs = append(teamIDs, team.TeamID)
	}

	t := query.Use(dal.DB()).Team
	return t.WithContext(ctx).Where(t.ID.In(teamIDs...)).Find()
}

func ListMembersByTeamID(ctx context.Context, teamID int64) ([]*model.User, error) {
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID)).Find()
	if err != nil {
		return nil, err
	}

	var userIDs []int64
	for _, team := range userTeams {
		userIDs = append(userIDs, team.UserID)
	}

	u := query.Use(dal.DB()).User
	return u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
}

func InviteMember(ctx context.Context, teamID, userID int64) error {
	return query.Use(dal.DB()).UserTeam.WithContext(ctx).Create(&model.UserTeam{
		UserID: userID,
		TeamID: teamID,
	})
}
