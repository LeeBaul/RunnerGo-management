package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
)

func SignUp(ctx context.Context, email, password, nickname string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{Email: email, Password: string(hashedPassword), Nickname: nickname}
	team := model.Team{Name: fmt.Sprintf("%s 的团队", nickname)}

	err = query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if err := tx.User.WithContext(ctx).Create(&user); err != nil {
			return err
		}

		if err := tx.Team.WithContext(ctx).Create(&team); err != nil {
			return err
		}

		return tx.UserTeam.WithContext(ctx).Create(&model.UserTeam{
			UserID: user.ID,
			TeamID: team.ID,
			RoleID: consts.RoleTypeOwner,
		})
	})

	if err != nil {
		return nil, err
	}

	SetUserSettings(ctx, user.ID, &rao.UserSettings{CurrentTeamID: team.ID})

	return &user, nil
}

func Login(ctx context.Context, email, password string) (*model.User, error) {
	tx := query.Use(dal.DB()).User
	user, err := tx.WithContext(ctx).Where(tx.Email.Eq(email)).First()
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
