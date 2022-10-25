package team

import (
	"context"
	"fmt"

	"github.com/go-omnibus/omnibus"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/mail"
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
			tx.Team.CreatedUserID.Value(userID),
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
	userTeams, err := ut.WithContext(ctx).Where(ut.UserID.Eq(userID)).Distinct(ut.TeamID).Find()
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

	var userIDs []int64
	for _, team := range teams {
		userIDs = append(userIDs, team.CreatedUserID)
	}
	u := dal.GetQuery().User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransTeamsModelToRaoTeam(teams, userTeams, teamCnt, users), nil
}

func ListMembersByTeamID(ctx context.Context, teamID int64) ([]*rao.Member, error) {
	ut := query.Use(dal.DB()).UserTeam
	userTeams, err := ut.WithContext(ctx).Where(ut.TeamID.Eq(teamID)).Order(ut.RoleID).Find()
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

func InviteMember(ctx context.Context, inviteUserID, teamID int64, members []*rao.InviteMember) (*rao.InviteMemberResp, error) {

	var emails []string
	memo := make(map[string]int64)
	for _, member := range members {
		emails = append(emails, member.Email)
		memo[member.Email] = member.RoleID
	}

	tx := query.Use(dal.DB()).User
	users, err := tx.WithContext(ctx).Where(tx.Email.In(emails...)).Find()
	if err != nil {
		return nil, err
	}

	var registerEmail []string
	for _, user := range users {
		registerEmail = append(registerEmail, user.Email)
	}
	registerEmail = omnibus.StringArrayUnique(registerEmail)

	var userIDs []int64
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	utt := dal.GetQuery().UserTeam
	existUser, err := utt.WithContext(ctx).Where(utt.TeamID.Eq(teamID), utt.UserID.In(userIDs...)).Find()
	if err != nil {
		return nil, err
	}

	for i, user := range users {
		for _, eu := range existUser {
			if eu.UserID == user.ID {
				users[i] = nil
			}
		}
	}

	var ut []*model.UserTeam
	for _, user := range users {
		if user != nil {
			ut = append(ut, &model.UserTeam{
				UserID:       user.ID,
				TeamID:       teamID,
				InviteUserID: inviteUserID,
				RoleID:       memo[user.Email],
			})
		}
	}

	if err := query.Use(dal.DB()).UserTeam.WithContext(ctx).CreateInBatches(ut, 5); err != nil {
		return nil, err
	}

	u, err := tx.WithContext(ctx).Where(tx.ID.Eq(inviteUserID)).First()
	if err != nil {
		return nil, err
	}

	px := dal.GetQuery().Team
	t, err := px.WithContext(ctx).Where(px.ID.Eq(teamID)).First()
	if err != nil {
		return nil, err
	}

	for _, e := range registerEmail {
		if err := mail.SendInviteEmail(ctx, e, u.Nickname, t.Name, true); err != nil {
			return nil, err
		}
	}

	unRegisterEmail := omnibus.StringArrayUnique(omnibus.StringArrayDiff(emails, registerEmail))
	if len(unRegisterEmail) > 0 {
		var userQueue []*model.TeamUserQueue
		for _, e := range unRegisterEmail {
			if err := mail.SendInviteEmail(ctx, e, u.Nickname, t.Name, false); err != nil {
				return nil, err
			}

			userQueue = append(userQueue, &model.TeamUserQueue{
				Email:  e,
				TeamID: teamID,
			})
		}
		qx := dal.GetQuery().TeamUserQueue
		if err := qx.WithContext(ctx).CreateInBatches(userQueue, 5); err != nil {
			return nil, err
		}
	}

	return &rao.InviteMemberResp{
		RegisterNum:      len(registerEmail),
		UnRegisterNum:    len(unRegisterEmail),
		UnRegisterEmails: unRegisterEmail,
	}, nil
}

func RoleUser(ctx context.Context, teamID, userID, roleID int64) error {
	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		_, err := tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(userID)).First()
		if err != nil {
			return err
		}

		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(teamID), tx.UserTeam.UserID.Eq(userID)).UpdateColumn(tx.UserTeam.RoleID, roleID)

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

		if team.CreatedUserID == userID {
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

func DisbandTeam(ctx context.Context, teamID, userID int64) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		t, err := tx.Team.WithContext(ctx).Where(
			tx.Team.ID.Eq(teamID), tx.Team.Type.Eq(consts.TeamTypeNormal), tx.Team.CreatedUserID.Eq(userID)).First()
		if err != nil {
			return err
		}

		_, err = tx.UserTeam.WithContext(ctx).Where(tx.UserTeam.TeamID.Eq(t.ID)).Delete()
		if err != nil {
			return err
		}

		settings, err := tx.Setting.WithContext(ctx).Where(tx.Setting.TeamID.Eq(teamID)).Find()
		if err != nil {
			return err
		}

		for _, s := range settings {
			pt, err := tx.Team.WithContext(ctx).Where(tx.Team.CreatedUserID.Eq(s.UserID), tx.Team.Type.Eq(consts.TeamTypePrivate)).First()
			if err != nil {
				return err
			}

			_, err = tx.Setting.WithContext(ctx).Where(tx.Setting.ID.Eq(s.ID)).UpdateColumn(tx.Setting.TeamID, pt.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func TransferTeam(ctx context.Context, teamID, userID, toUserID int64) error {
	return nil
}
