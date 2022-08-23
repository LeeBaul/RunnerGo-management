package packer

import (
	"fmt"

	"github.com/bytedance/sonic"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransGroupReqToGroup(group *rao.SaveGroupReq) *mao.Group {
	reqByte, err := sonic.Marshal(group.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("group.request json marshal err %w", err))
	}

	scriptByte, err := sonic.Marshal(group.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("group.script json marshal err %w", err))
	}

	return &mao.Group{
		TargetID: group.TargetID,
		Request:  reqByte,
		Script:   scriptByte,
	}
}
