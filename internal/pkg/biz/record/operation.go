package record

import (
	"context"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
)

func InsertCreate(ctx context.Context, teamID, userID int64, operate int32, name string) error {
	return insert(ctx, teamID, userID, name, consts.OperationCategoryCreate, operate)
}

func InsertUpdate(ctx context.Context, teamID, userID int64, operate int32, name string) error {
	return insert(ctx, teamID, userID, name, consts.OperationCategoryUpdate, operate)
}

func InsertDelete(ctx context.Context, teamID, userID int64, operate int32, name string) error {
	return insert(ctx, teamID, userID, name, consts.OperationCategoryDelete, operate)
}

func InsertRun(ctx context.Context, teamID, userID int64, operate int32, name string) error {
	return insert(ctx, teamID, userID, name, consts.OperationCategoryRun, operate)
}

func insert(ctx context.Context, teamID, userID int64, name string, category, operate int32) error {
	return query.Use(dal.DB()).Operation.WithContext(ctx).Create(&model.Operation{
		TeamID:   teamID,
		UserID:   userID,
		Category: category,
		Name:     name,
		Operate:  operate,
	})
}
