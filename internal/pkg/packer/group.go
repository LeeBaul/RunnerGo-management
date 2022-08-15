package packer

import (
	"encoding/json"
	"fmt"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransGroupReqToGroup(group *rao.SaveGroupReq) *model.Group {
	reqByte, err := json.Marshal(group.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("group.request json marshal err %w", err))
	}

	scriptByte, err := json.Marshal(group.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("group.script json marshal err %w", err))
	}

	return &model.Group{
		TargetID: group.TargetID,
		Request:  string(reqByte),
		Script:   string(scriptByte),
	}
}
