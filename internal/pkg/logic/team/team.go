package team

import (
	"context"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func ListByUserID(ctx context.Context, userID int64) ([]*rao.Team, error) {
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
	teams, err := t.WithContext(ctx).Where(t.ID.In(teamIDs...)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransTeamsModelToResp(teams, userTeams), nil
}

func ListMembersByTeamID(ctx context.Context, teamID int64) ([]*rao.Member, error) {
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
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransUsersModelToMembers(users, userTeams), nil
}

func InviteMember(ctx context.Context, teamID int64, email []string) error {
	tx := query.Use(dal.DB()).User
	users, err := tx.WithContext(ctx).Where(tx.Email.In(email...)).Find()
	if err != nil {
		return err
	}

	var ut []*model.UserTeam
	for _, u := range users {
		ut = append(ut, &model.UserTeam{
			UserID: u.ID,
			TeamID: teamID,
			RoleID: consts.RoleTypeMember,
		})
	}

	return query.Use(dal.DB()).UserTeam.WithContext(ctx).CreateInBatches(ut, 5)
}

func RemoveMember(ctx context.Context, teamID, memberID int64) error {
	tx := query.Use(dal.DB()).UserTeam
	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.UserID.Eq(memberID)).Delete()
	return err
}
