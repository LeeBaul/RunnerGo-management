package machine

import (
	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/proof"
	"gorm.io/gen"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/rao"
)

func GetMachineList(ctx *gin.Context, req rao.GetMachineListParam) ([]*rao.GetMachineListResponse, error) {
	// 查询机器列表
	tx := dal.GetQuery().Machine

	conditions := make([]gen.Condition, 0)
	if req.Name != "" {
		conditions = append(conditions, tx.Name.Eq(req.Name))
	}
	if req.ServerType != 0 {
		conditions = append(conditions, tx.ServerType.Eq(req.ServerType))
	}

	machineList, err := tx.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		proof.Errorf("机器列表--获取机器列表数据失败，err:", err)
		return nil, err
	}

	res := make([]*rao.GetMachineListResponse, 0, len(machineList))

	for _, machineInfo := range machineList {
		machineTmp := &rao.GetMachineListResponse{
			Region:     machineInfo.Region,
			IP:         machineInfo.IP,
			Port:       machineInfo.Port,
			Name:       machineInfo.Name,
			ServerType: machineInfo.ServerType,
			Status:     machineInfo.Status,
		}
		res = append(res, machineTmp)
	}

	return res, nil
}
