package target

import (
	"context"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func ListFolderAPI(ctx context.Context, teamID int64, limit, offset int) ([]*rao.FolderAPI, int64, error) {

	tx := query.Use(dal.DB()).Target
	targets, err := tx.WithContext(ctx).
		Where(tx.TeamID.Eq(teamID), tx.TargetType.In(consts.TargetTypeFolder, consts.TargetTypeAPI)).
		Order(tx.Sort.Desc(), tx.CreatedAt.Desc()).
		Limit(limit).Offset(offset).Find()

	if err != nil {
		return nil, 0, err
	}

	cnt, err := tx.WithContext(ctx).
		Where(tx.TeamID.Eq(teamID), tx.TargetType.In(consts.TargetTypeFolder, consts.TargetTypeAPI)).
		Count()

	if err != nil {
		return nil, 0, err
	}

	return packer.TransTargetToFolderAPI(targets), cnt, nil
}
