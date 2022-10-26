package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	services "kp-management/api"
	"time"

	//services "kp-controller/api"
	"kp-management/internal/pkg/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/logic/stress"
)

//func RunStress(ctx context.Context, req *services.RunStressReq) (*services.RunStressResp, error) {
func RunStress(ctx *gin.Context, runStressParams RunStressParams) error {
	rms := &stress.RunMachineStress{}

	siv := &stress.SplitImportVariable{}
	siv.SetNext(rms)

	ss := &stress.SplitStress{}
	ss.SetNext(siv)

	ms := &stress.MakeStress{}
	ms.SetNext(ss)

	mr := &stress.MakeReport{}
	mr.SetNext(ms)

	iv := &stress.AssembleImportVariables{}
	iv.SetNext(mr)

	sv := &stress.AssembleSceneVariables{}
	sv.SetNext(iv)

	f := &stress.AssembleFlows{}
	f.SetNext(sv)

	v := &stress.AssembleGlobalVariables{}
	v.SetNext(f)

	t := &stress.AssembleTask{}
	t.SetNext(v)

	s := &stress.AssembleScenes{}
	s.SetNext(t)

	p := &stress.AssemblePlan{}
	p.SetNext(s)

	m := &stress.CheckIdleMachine{}
	m.SetNext(p)

	err := m.Execute(&stress.Baton{
		Ctx:      ctx,
		PlanID:   runStressParams.PlanID,
		TeamID:   runStressParams.TeamID,
		SceneIDs: runStressParams.SceneID,
		UserID:   runStressParams.UserID,
	})

	return err
}

func NotifyStopStress(ctx context.Context, req *services.NotifyStopStressReq) (*services.NotifyStopStressResp, error) {

	err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		r := tx.Report
		// 修改报告状态
		_, err := r.WithContext(ctx).Where(r.ID.Eq(req.GetReportID())).UpdateSimple(r.Status.Value(consts.ReportStatusFinish), r.UpdatedAt.Value(time.Now()))
		if err != nil {
			return err
		}

		// 查找报告对应计划
		report, err := r.WithContext(ctx).Where(r.ID.Eq(req.GetReportID())).First()
		if err != nil {
			return err
		}

		// 统计报告是否全部完成
		reportCnt, err := r.WithContext(ctx).Where(r.PlanID.Eq(report.PlanID)).Count()
		if err != nil {
			return err
		}
		finishReportCnt, err := r.WithContext(ctx).Where(r.PlanID.Eq(report.PlanID), r.Status.Eq(consts.ReportStatusFinish)).Count()
		if err != nil {
			return err
		}

		if finishReportCnt == reportCnt { // 报告全部完成则计划也完成
			p := tx.Plan
			_, err := p.WithContext(ctx).Where(p.ID.Eq(report.PlanID)).UpdateSimple(p.Status.Value(consts.PlanStatusNormal), p.UpdatedAt.Value(time.Now()))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &services.NotifyStopStressResp{
		Code: 0,
		Msg:  "ok",
	}, nil
}
