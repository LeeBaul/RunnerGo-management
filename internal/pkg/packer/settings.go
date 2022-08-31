package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransUserSettingsToRaoUserSettings(s *model.Setting) *rao.UserSettings {
	return &rao.UserSettings{
		CurrentTeamID: s.TeamID,
	}
}
