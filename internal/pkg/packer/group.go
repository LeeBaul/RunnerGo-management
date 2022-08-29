package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransGroupReqToGroup(group *rao.SaveGroupReq) *mao.Group {
	request, err := bson.Marshal(group.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("group.request json marshal err %w", err))
	}

	script, err := bson.Marshal(group.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("group.script json marshal err %w", err))
	}

	return &mao.Group{
		TargetID: group.TargetID,
		Request:  request,
		Script:   script,
	}
}
