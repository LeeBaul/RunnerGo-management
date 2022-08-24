package auth

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func SetUserSettings(ctx context.Context, userID int64, settings *rao.UserSettings) error {
	tx := query.Use(dal.DB()).Setting
	return tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).Save(&model.Setting{
		UserID: userID,
		TeamID: settings.CurrentTeamID,
	})
}

func GetUserSettings(ctx context.Context, userID int64) (*rao.UserSettings, error) {
	tx := query.Use(dal.DB()).Setting
	s, err := tx.WithContext(ctx).Where(tx.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	return packer.TransUserSettingsToResp(s), nil
}
