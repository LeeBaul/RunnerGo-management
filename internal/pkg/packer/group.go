package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
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

func TransTargetToGroup(t *model.Target, g *mao.Group) *rao.Group {
	var r rao.Request
	if err := bson.Unmarshal(g.Request, &r); err != nil {
		fmt.Sprintln(fmt.Errorf("folder.request json UnMarshal err %w", err))
	}

	var s rao.Script
	if err := bson.Unmarshal(g.Script, &s); err != nil {
		fmt.Sprintln(fmt.Errorf("folder.script json UnMarshal err %w", err))
	}

	return &rao.Group{
		TeamID:   t.TeamID,
		TargetID: t.ID,
		ParentID: t.ParentID,
		Name:     t.Name,
		Method:   t.Method,
		Sort:     t.Sort,
		TypeSort: t.TypeSort,
		Version:  t.Version,
		Request:  &r,
		Script:   &s,
	}
}
