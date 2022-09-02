package folder

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, userID int64, req *rao.SaveFolderReq) error {
	target := packer.TransSaveFolderReqToTargetModel(req, userID)
	folder := packer.TransSaveFolderReqToMaoFolder(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFolder)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			folder.TargetID = target.ID
			_, err := collection.InsertOne(ctx, folder)
			if err != nil {
				return err
			}

			return record.InsertCreate(ctx, target.TeamID, userID, fmt.Sprintf("创建文件夹 - %s", target.Name))
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": folder})
		if err != nil {
			return err
		}

		return record.InsertUpdate(ctx, target.TeamID, userID, fmt.Sprintf("修改文件夹 - %s", target.Name))
	})
}

func GetByTargetID(ctx context.Context, teamID, targetID int64) (*rao.Folder, error) {
	tx := query.Use(dal.DB()).Target
	t, err := tx.WithContext(ctx).Where(
		tx.ID.Eq(targetID),
		tx.TeamID.Eq(teamID),
		tx.TargetType.Eq(consts.TargetTypeFolder),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
	).First()

	if err != nil {
		return nil, err
	}

	var f *mao.Folder
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFolder)
	err = collection.FindOne(ctx, bson.D{{"target_id", targetID}}).Decode(&f)
	if err != nil {
		return nil, err
	}

	return packer.TransTargetToRaoFolder(t, f), nil
}
