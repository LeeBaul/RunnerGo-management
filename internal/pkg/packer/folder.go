package packer

import (
	"encoding/json"
	"fmt"
	"kp-management/internal/pkg/dal/mao"

	"kp-management/internal/pkg/dal/rao"
)

func TransFolderReqToFolder(folder *rao.SaveFolderReq) *mao.Folder {
	reqByte, err := json.Marshal(folder.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.request json marshal err %w", err))
	}

	scriptByte, err := json.Marshal(folder.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.script json marshal err %w", err))
	}

	//return &model.Folder{
	//	TargetID: folder.TargetID,
	//	Request:  string(reqByte),
	//	Script:   string(scriptByte),
	//}

	return &mao.Folder{
		TargetID: folder.TargetID,
		Request:  string(reqByte),
		Script:   string(scriptByte),
	}
}
