package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransFolderReqToFolder(folder *rao.SaveFolderReq) *mao.Folder {
	request, err := bson.Marshal(folder.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.request json marshal err %w", err))
	}

	script, err := bson.Marshal(folder.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.script json marshal err %w", err))
	}

	return &mao.Folder{
		TargetID: folder.TargetID,
		Request:  request,
		Script:   script,
	}
}
