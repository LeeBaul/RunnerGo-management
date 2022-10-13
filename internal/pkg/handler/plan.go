package handler

import (
	"github.com/go-omnibus/omnibus"
	"github.com/go-resty/resty/v2"

	services "kp-management/api"
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/dal/runner"
	"kp-management/internal/pkg/logic/plan"

	"github.com/gin-gonic/gin"
)

func RunPlan(ctx *gin.Context) {
	var req rao.RunPlanReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	_, err := dal.ClientGRPC().RunStress(ctx, &services.RunStressReq{
		PlanID:  req.PlanID,
		TeamID:  req.TeamID,
		SceneID: req.SceneID,
		UserID:  jwt.GetUserIDByCtx(ctx),
	})

	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrHttpFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func StopPlan(ctx *gin.Context) {
	var req rao.StopPlanReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().Report
	reports, err := tx.WithContext(ctx).Where(tx.PlanID.In(req.PlanIDs...)).Find()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	var reportIDs []int64
	for _, report := range reports {
		reportIDs = append(reportIDs, report.ID)
	}
	_, err = resty.New().R().
		SetBody(runner.StopRunnerReq{ReportIds: omnibus.Int64sToStrings(reportIDs)}).
		Post(conf.Conf.Clients.Runner.StopPlan)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrHttpFailed, err.Error())
		return
	}

	px := dal.GetQuery().Plan
	_, err = px.WithContext(ctx).Where(px.ID.In(req.PlanIDs...)).UpdateColumn(px.Status, consts.PlanStatusNormal)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	//_, err = tx.WithContext(ctx).Where(tx.PlanID.In(req.PlanIDs...)).UpdateColumn(tx.Status, consts.ReportStatusFinish)
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
	//	return
	//}

	response.Success(ctx)
	return
}

// ClonePlan 克隆计划
func ClonePlan(ctx *gin.Context) {
	var req rao.ClonePlanReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := plan.ClonePlan(ctx, req.PlanID, jwt.GetUserIDByCtx(ctx)); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// ListUnderwayPlan 运行中的计划
func ListUnderwayPlan(ctx *gin.Context) {
	var req rao.ListUnderwayPlanReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	plans, total, err := plan.ListByStatus(ctx, req.TeamID, consts.PlanStatusUnderway, req.Size, (req.Page-1)*req.Size)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListUnderwayPlanResp{
		Plans: plans,
		Total: total,
	})
	return
}

// ListPlans 测试计划列表
func ListPlans(ctx *gin.Context) {
	var req rao.ListPlansReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	plans, total, err := plan.ListByTeamID(ctx, req.TeamID, req.Size, (req.Page-1)*req.Size, req.Keyword, req.StartTimeSec, req.EndTimeSec)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListPlansResp{
		Plans: plans,
		Total: total,
	})
	return
}

// SavePlan 创建修改计划
func SavePlan(ctx *gin.Context) {
	var req rao.SavePlanReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	planID, err := plan.Save(ctx, &req, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SavePlanResp{PlanID: planID})
	return
}

// SavePlanTask 创建/修改计划配置
func SavePlanTask(ctx *gin.Context) {
	var req rao.SavePlanConfReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := plan.SaveTask(ctx, &req, jwt.GetUserIDByCtx(ctx)); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// GetPlan 获取计划
func GetPlan(ctx *gin.Context) {
	var req rao.GetPlanConfReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	p, err := plan.GetByPlanID(ctx, req.TeamID, req.PlanID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetPlanResp{Plan: p})
	return
}

// DeletePlan 删除计划
func DeletePlan(ctx *gin.Context) {
	var req rao.DeletePlanReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := plan.DeleteByPlanID(ctx, req.TeamID, req.PlanID); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// SetPreinstall 保存预设设置
func SetPreinstall(ctx *gin.Context) {
	var req rao.SetPreinstallReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := plan.SetPreinstall(ctx, &req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// GetPreinstall 获取预设设置
func GetPreinstall(ctx *gin.Context) {
	var req rao.GetPreinstallReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	p, err := plan.GetPreinstall(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetPreinstallResp{Preinstall: p})
}
