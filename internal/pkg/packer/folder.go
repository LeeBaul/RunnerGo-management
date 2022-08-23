package packer

import (
	"fmt"

	"github.com/bytedance/sonic"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransFolderReqToFolder(folder *rao.SaveFolderReq) *mao.Folder {
	request, err := sonic.MarshalString(folder.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.request json marshal err %w", err))
	}

	script, err := sonic.MarshalString(folder.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.script json marshal err %w", err))
	}

	return &mao.Folder{
		TargetID: folder.TargetID,
		Request:  request,
		Script:   script,
	}
}
