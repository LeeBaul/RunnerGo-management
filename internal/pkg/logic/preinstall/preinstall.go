package preinstall

import (
	"github.com/go-omnibus/proof"
	"github.com/goccy/go-json"
	"golang.org/x/net/context"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func SavePreinstall(ctx context.Context, req *rao.SavePreinstallReq) (int, error) {
	tx := dal.GetQuery().PreinstallConf

	// 把mode_conf压缩成字符串
	modeConfString, err := json.Marshal(req.ModeConf)
	if err != nil {
		proof.Errorf("保存预设配置--压缩mode_conf为字符串失败，err:", err)
		return errno.ErrMarshalFailed, err
	}
	// 把timed_task_conf压缩成字符串
	timedTaskConfString, err := json.Marshal(req.TimedTaskConf)
	if err != nil {
		proof.Errorf("保存预设配置--压缩timed_task_conf为字符串失败，err:", err)
		return errno.ErrMarshalFailed, err
	}

	if req.ID == 0 { // 新建
		// 排重
		_, err = tx.WithContext(ctx).Where(tx.ConfName.Eq(req.ConfName)).First()
		if err == nil {
			proof.Infof("保存预设配置--查询预设配置表失败,或已存在，err:", err)
			return errno.ErrMysqlFailed, err
		}

		insertData := &model.PreinstallConf{
			ConfName:      req.ConfName,
			UserID:        req.UserID,
			UserName:      req.UserName,
			TaskType:      req.TaskType,
			TaskMode:      req.TaskMode,
			ModeConf:      string(modeConfString),
			TimedTaskConf: string(timedTaskConfString),
		}
		err = tx.WithContext(ctx).Create(insertData)
		if err != nil {
			proof.Errorf("保存预设配置--创建数据失败，err:", err)
			return errno.ErrMysqlFailed, err
		}
	} else { // 修改
		updateData := model.PreinstallConf{
			ConfName:      req.ConfName,
			UserID:        req.UserID,
			UserName:      req.UserName,
			TaskType:      req.TaskType,
			TaskMode:      req.TaskMode,
			ModeConf:      string(modeConfString),
			TimedTaskConf: string(timedTaskConfString),
		}
		_, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.ID)).Updates(updateData)
		if err != nil {
			proof.Errorf("保存预设配置--修改数据失败，err:", err)
			return errno.ErrMysqlFailed, err
		}
	}
	return errno.Ok, nil
}
