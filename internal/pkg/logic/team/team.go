package team

import (
	"context"
	"fmt"

	"github.com/go-omnibus/omnibus"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func SaveTeam(ctx context.Context, teamID, userID int64, name string) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		team, err := tx.Team.WithContext(ctx).Where(tx.Team.ID.Eq(teamID)).Assign(
			tx.Team.ID.Value(teamID),
			tx.Team.Name.Value(name),
			tx.Team.Type.Value(consts.TeamTypeNormal),
		).FirstOrCreate()
		if err != nil {
			return err
		}

		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(team.ID), tx.UserTeam.UserID.Eq(userID)).Assign(
			tx.UserTeam.TeamID.Value(team.ID),
			tx.UserTeam.UserID.Value(userID),
			tx.UserTeam.RoleID.Value(consts.RoleTypeOwner),
		).FirstOrCreate()

		return err
	})
}

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

	var teamCnt []*packer.TeamMemberCount
	if err := ut.WithContext(ctx).Select(ut.TeamID, ut.UserID.Count().As("cnt")).Where(ut.TeamID.In(teamIDs...)).Group(ut.TeamID).Scan(&teamCnt); err != nil {
		return nil, err
	}

	return packer.TransTeamsModelToRaoTeam(teams, userTeams, teamCnt), nil
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
		userIDs = append(userIDs, team.InviteUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(omnibus.Int64ArrayUnique(userIDs)...)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransUsersToRaoMembers(users, userTeams), nil
}

func InviteMember(ctx context.Context, inviteUserID, teamID int64, email []string) error {
	var userIDs []int64

	tx := query.Use(dal.DB()).User
	err := tx.WithContext(ctx).Where(tx.Email.In(email...)).Pluck(tx.ID, &userIDs)
	if err != nil {
		return err
	}

	var ut []*model.UserTeam
	for _, userID := range userIDs {
		ut = append(ut, &model.UserTeam{
			UserID:       userID,
			TeamID:       teamID,
			InviteUserID: inviteUserID,
			RoleID:       consts.RoleTypeMember,
		})
	}

	return query.Use(dal.DB()).UserTeam.WithContext(ctx).CreateInBatches(ut, 5)
}

func RoleUser(ctx context.Context, teamID, userID, roleID int64) error {
	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		u, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(userID)).First()
		if err != nil {
			return err
		}

		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.ID.Eq(u.ID)).UpdateColumn(tx.UserTeam.RoleID, roleID)

		return err
	})
}

func RemoveMember(ctx context.Context, teamID, userID, memberID int64) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		// 不能移除自己
		if userID == memberID {
			return fmt.Errorf("user no permissions")
		}

		admin, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(userID)).First()
		if err != nil {
			return err
		}

		// 只有管理员和创建人可以操作移除
		if !omnibus.InArray(admin.RoleID, []int64{consts.RoleTypeAdmin, consts.RoleTypeOwner}) {
			return fmt.Errorf("user no permissions")
		}

		user, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(memberID)).First()
		if err != nil {
			return err
		}

		// 不能移除创建人
		if user.RoleID == consts.RoleTypeOwner {
			return fmt.Errorf("user no permissions")
		}

		// 只有创建人能移除管理员
		if user.RoleID == consts.RoleTypeAdmin && admin.RoleID != consts.RoleTypeOwner {
			return fmt.Errorf("user no permissions")
		}

		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(memberID)).Delete()

		return err
	})
}

func QuitTeam(ctx context.Context, teamID, userID int64) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		team, err := tx.Team.WithContext(ctx).Where(tx.Team.ID.Eq(teamID)).First()
		if err != nil {
			return err
		}

		if team.Type == consts.TeamTypePrivate {
			return fmt.Errorf("user no permissions")
		}

		ut, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(userID)).First()
		if err != nil {
			return err
		}

		switch ut.RoleID {
		case consts.RoleTypeOwner:
			cnt, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.RoleID.Eq(consts.RoleTypeAdmin)).Count()
			if err != nil {
				return err
			}
			if cnt == 0 {
				return fmt.Errorf("not found admin user")
			}
		case consts.RoleTypeMember, consts.RoleTypeAdmin:
			break
		}

		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(userID)).Delete()
		if err != nil {
			return err
		}

		return nil
	})
}
