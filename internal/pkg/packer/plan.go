package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransPlansToResp(plans []*model.Plan) []*rao.Plan {
	ret := make([]*rao.Plan, 0)
	for _, plan := range plans {
		ret = append(ret, &rao.Plan{
			PlanID:     plan.ID,
			Name:       plan.Name,
			UpdatedSec: plan.UpdatedAt.Unix(),
		})
	}

	return ret
}
