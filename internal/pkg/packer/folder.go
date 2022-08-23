package packer

import (
	"fmt"

	"github.com/bytedance/sonic"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransFolderReqToFolder(folder *rao.SaveFolderReq) *mao.Folder {
	reqByte, err := sonic.Marshal(folder.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.request json marshal err %w", err))
	}

	scriptByte, err := sonic.Marshal(folder.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.script json marshal err %w", err))
	}

	return &mao.Folder{
		TargetID: folder.TargetID,
		Request:  reqByte,
		Script:   scriptByte,
	}
}
