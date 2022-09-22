package report

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"kp-management/internal/pkg/conf"

	"github.com/go-omnibus/proof"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gen"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func CountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Report

	return tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Count()
}

func ListByTeamID2(ctx context.Context, teamID int64, limit, offset int, keyword string, startTimeSec, endTimeSec int64) ([]*rao.Report, int64, error) {

	tx := query.Use(dal.DB()).Report

	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.TeamID.Eq(teamID))

	if keyword != "" {
		var reportIDs []int64

		planReportIDs, err := KeywordFindPlan(ctx, teamID, keyword)
		if err != nil {
			return nil, 0, err
		}
		reportIDs = append(reportIDs, planReportIDs...)

		sceneReportIDs, err := KeywordFindScene(ctx, teamID, keyword)
		if err != nil {
			return nil, 0, err
		}
		reportIDs = append(reportIDs, sceneReportIDs...)

		userReportIDs, err := KeywordFindUser(ctx, keyword)
		if err != nil {
			return nil, 0, err
		}
		reportIDs = append(reportIDs, userReportIDs...)

		if len(reportIDs) > 0 {
			conditions = append(conditions, tx.ID.In(reportIDs...))
		} else {
			conditions = append(conditions, tx.ID.In(0))
		}
	}

	if startTimeSec > 0 && endTimeSec > 0 {
		startTime := time.Unix(startTimeSec, 0)
		endTime := time.Unix(endTimeSec, 0)
		conditions = append(conditions, tx.CreatedAt.Between(startTime, endTime))
	}

	reports, cnt, err := tx.WithContext(ctx).Where(conditions...).
		Order(tx.ID.Desc()).
		FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range reports {
		userIDs = append(userIDs, r.RunUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	var planIDs []int64
	var sceneIDs []int64
	for _, report := range reports {
		planIDs = append(planIDs, report.PlanID)
		sceneIDs = append(sceneIDs, report.SceneID)
	}

	p := dal.GetQuery().Plan
	plans, err := p.WithContext(ctx).Where(p.ID.In(planIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	s := dal.GetQuery().Target
	scenes, err := s.WithContext(ctx).Where(s.ID.In(sceneIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransReportModelToRaoReportList(reports, users, plans, scenes), cnt, nil
}

func KeywordFindPlan(ctx context.Context, teamID int64, keyword string) ([]int64, error) {
	var planIDs []int64

	p := dal.GetQuery().Plan
	err := p.WithContext(ctx).Where(p.TeamID.Eq(teamID), p.Name.Like(fmt.Sprintf("%%%s%%", keyword))).Pluck(p.ID, &planIDs)
	if err != nil {
		return nil, err
	}

	if len(planIDs) == 0 {
		return nil, nil
	}

	var reportIDs []int64
	r := dal.GetQuery().Report
	err = r.WithContext(ctx).Where(r.PlanID.In(planIDs...)).Pluck(r.ID, &reportIDs)
	if err != nil {
		return nil, err
	}

	return reportIDs, nil
}

func KeywordFindScene(ctx context.Context, teamID int64, keyword string) ([]int64, error) {
	var sceneIDs []int64

	s := dal.GetQuery().Target
	err := s.WithContext(ctx).Where(s.TeamID.Eq(teamID), s.Name.Like(fmt.Sprintf("%%%s%%", keyword))).Pluck(s.ID, &sceneIDs)
	if err != nil {
		return nil, err
	}

	if len(sceneIDs) == 0 {
		return nil, nil
	}

	var reportIDs []int64
	r := dal.GetQuery().Report
	err = r.WithContext(ctx).Where(r.SceneID.In(sceneIDs...)).Pluck(r.ID, &reportIDs)
	if err != nil {
		return nil, err
	}

	return reportIDs, nil
}

func KeywordFindUser(ctx context.Context, keyword string) ([]int64, error) {
	var userIDs []int64

	u := query.Use(dal.DB()).User
	err := u.WithContext(ctx).Where(u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword))).Pluck(u.ID, &userIDs)
	if err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return nil, nil
	}

	var reportIDs []int64
	r := dal.GetQuery().Report
	err = r.WithContext(ctx).Where(r.RunUserID.In(userIDs...)).Pluck(r.ID, &reportIDs)
	if err != nil {
		return nil, err
	}

	return reportIDs, nil
}

func ListByTeamID(ctx context.Context, teamID int64, limit, offset int, keyword string, startTimeSec, endTimeSec int64) ([]*rao.Report, int64, error) {
	tx := query.Use(dal.DB()).Report

	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.TeamID.Eq(teamID))

	if keyword != "" {
		p := dal.GetQuery().Plan
		plans, err := p.WithContext(ctx).Where(p.Name.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		if err != nil {
			return nil, 0, err
		}
		var planIDs []int64
		for _, plan := range plans {
			planIDs = append(planIDs, plan.ID)
		}
		conditions = append(conditions, tx.PlanID.In(planIDs...))

		//s := dal.GetQuery().Target
		//scenes, err := s.WithContext(ctx).Where(s.Name.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		//if err != nil {
		//	return nil, 0, err
		//}
		//var sceneIDs []int64
		//for _, scene := range scenes {
		//	sceneIDs = append(sceneIDs, scene.ID)
		//}
		//if len(sceneIDs) > 0 {
		//	conditions[1] = tx.SceneID.In(sceneIDs...)
		//}

		u := query.Use(dal.DB()).User
		users, err := u.WithContext(ctx).Where(u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		if err != nil {
			return nil, 0, err
		}

		if len(users) > 0 {
			conditions[1] = tx.RunUserID.Eq(users[0].ID)
		}
	}

	if startTimeSec > 0 && endTimeSec > 0 {
		startTime := time.Unix(startTimeSec, 0)
		endTime := time.Unix(endTimeSec, 0)
		conditions = append(conditions, tx.CreatedAt.Between(startTime, endTime))
	}

	reports, cnt, err := tx.WithContext(ctx).Where(conditions...).
		Order(tx.ID.Desc()).
		FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range reports {
		userIDs = append(userIDs, r.RunUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	var planIDs []int64
	var sceneIDs []int64
	for _, report := range reports {
		planIDs = append(planIDs, report.PlanID)
		sceneIDs = append(sceneIDs, report.SceneID)
	}

	p := dal.GetQuery().Plan
	plans, err := p.WithContext(ctx).Where(p.ID.In(planIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	s := dal.GetQuery().Target
	scenes, err := s.WithContext(ctx).Where(s.ID.In(sceneIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransReportModelToRaoReportList(reports, users, plans, scenes), cnt, nil
}

func DeleteReport(ctx context.Context, teamID, reportID int64) error {
	tx := query.Use(dal.DB()).Report
	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(reportID)).Delete()

	return err
}

func GetTaskDetail(ctx context.Context, req rao.GetReportReq) (*rao.ReportTask, error) {
	var detail mao.ReportTask
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectReportTask)

	err := collection.FindOne(ctx, bson.D{{"report_id", req.ReportID}}).Decode(&detail)
	if err != nil {
		proof.Error("mongo decode err", proof.WithError(err))
		return nil, err
	}

	r := query.Use(dal.DB()).Report
	ru, err := r.WithContext(ctx).Where(r.TeamID.Eq(req.TeamID), r.ID.Eq(req.ReportID)).First()
	if err != nil {
		proof.Error("req not found err", proof.WithError(err))
		return nil, err
	}

	u := query.Use(dal.DB()).User
	user, err := u.WithContext(ctx).Where(u.ID.Eq(ru.RunUserID)).First()
	if err != nil {
		proof.Error("user not found err", proof.WithError(err))
		return nil, err
	}

	return &rao.ReportTask{
		UserID:         user.ID,
		UserName:       user.Nickname,
		UserAvatar:     user.Avatar,
		PlanID:         detail.PlanID,
		PlanName:       detail.PlanName,
		ReportID:       detail.ReportID,
		CreatedTimeSec: ru.CreatedAt.Unix(),
		TaskType:       detail.TaskType,
		TaskMode:       detail.TaskMode,
		TaskStatus:     ru.Status,
		ModeConf: &rao.ModeConf{
			ReheatTime:       detail.ModeConf.ReheatTime,
			RoundNum:         detail.ModeConf.RoundNum,
			Concurrency:      detail.ModeConf.Concurrency,
			ThresholdValue:   detail.ModeConf.ThresholdValue,
			StartConcurrency: detail.ModeConf.StartConcurrency,
			Step:             detail.ModeConf.Step,
			StepRunTime:      detail.ModeConf.StepRunTime,
			MaxConcurrency:   detail.ModeConf.MaxConcurrency,
			Duration:         detail.ModeConf.Duration,
		},
	}, nil
}

func GetReportDebugLog(ctx context.Context, report rao.GetReportReq) (err error, debugMsgList []map[string]interface{}) {
	//clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s", user, password, host, db))

	reportId := strconv.FormatInt(report.ReportID, 10)
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectStressDebug)
	filter := bson.D{{"report_id", reportId}}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		proof.Error("debug日志查询失败", proof.WithError(err))
		return
	}
	for cur.Next(ctx) {
		debugMsg := make(map[string]interface{})
		err = cur.Decode(&debugMsg)
		if err != nil {
			proof.Error("debug日志转换失败", proof.WithError(err))
			return
		}
		if debugMsg["end"] != true {
			debugMsgList = append(debugMsgList, debugMsg)
		}
	}
	return
}

func GetReportDetail(ctx context.Context, report rao.GetReportReq, host, user, password string) (err error, resultData ResultData) {
	reportId := strconv.FormatInt(report.ReportID, 10)
	//index := strconv.FormatInt(report.TeamID, 10)
	index := conf.Conf.ES.Index

	queryEs := elastic.NewBoolQuery()
	queryEs = queryEs.Must(elastic.NewMatchQuery("report_id", reportId))

	client, _ := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, password),
		elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
		elastic.SetHealthcheckInterval(30*time.Second),
	)
	_, _, err = client.Ping(host).Do(ctx)
	if err != nil {
		proof.Error("es连接失败", proof.WithError(err))
		return
	}
	//res, err := client.Search(index).Query(query).From(0).Size(size).Pretty(true).Do(context.Background())
	res, err := client.Search(index).Query(queryEs).Sort("time_stamp", true).Size(conf.Conf.ES.Size).Pretty(true).Do(ctx)
	if err != nil {
		proof.Error("获取报告详情失败", proof.WithError(err))
		return
	}
	if res == nil {
		proof.Error("报告详情为空")
		return
	}

	var resultMsg SceneTestResultDataMsg // 从es中获取得数据结构

	for _, item := range res.Hits.Hits {
		err = json.Unmarshal(item.Source, &resultMsg)
		if resultData.Results == nil {
			resultData.Results = make(map[string]*ResultDataMsg)
		}
		resultData.ReportId = resultMsg.ReportId
		resultData.End = resultMsg.End
		resultData.ReportName = resultMsg.ReportName
		resultData.PlanId = resultMsg.PlanId
		resultData.PlanName = resultMsg.PlanName
		resultData.SceneId = resultMsg.SceneId
		resultData.SceneName = resultMsg.SceneName
		resultData.TimeStamp = resultMsg.TimeStamp
		if resultMsg.Results != nil && len(resultMsg.Results) > 0 {
			for k, apiResult := range resultMsg.Results {
				if resultData.Results[k] == nil {
					resultData.Results[k] = new(ResultDataMsg)
				}
				resultData.Results[k].ApiName = apiResult.Name
				resultData.Results[k].Concurrency = apiResult.Concurrency
				resultData.Results[k].TotalRequestNum = apiResult.TotalRequestNum
				resultData.Results[k].TotalRequestTime = apiResult.TotalRequestTime
				resultData.Results[k].SuccessNum = apiResult.SuccessNum
				resultData.Results[k].ErrorNum = apiResult.ErrorNum
				if resultData.Results[k].ErrorNum != 0 {
					resultData.Results[k].ErrorRate = float64(apiResult.ErrorNum) / float64(apiResult.TotalRequestNum)
				}

				resultData.Results[k].AvgRequestTime = apiResult.AvgRequestTime
				resultData.Results[k].MaxRequestTime = apiResult.MaxRequestTime
				resultData.Results[k].MinRequestTime = apiResult.MinRequestTime
				resultData.Results[k].CustomRequestTimeLine = apiResult.CustomRequestTimeLine
				resultData.Results[k].CustomRequestTimeLineValue = apiResult.CustomRequestTimeLineValue
				resultData.Results[k].NinetyRequestTimeLine = apiResult.NinetyRequestTimeLine
				resultData.Results[k].NinetyFiveRequestTimeLine = apiResult.NinetyFiveRequestTimeLine
				resultData.Results[k].NinetyNineRequestTimeLine = apiResult.NinetyNineRequestTimeLine
				resultData.Results[k].SendBytes = apiResult.SendBytes
				resultData.Results[k].ReceivedBytes = apiResult.ReceivedBytes

				if resultData.Results[k].QpsList == nil {
					resultData.Results[k].QpsList = []TimeValue{}
				}
				var timeValue = TimeValue{}
				timeValue.TimeStamp = resultData.TimeStamp
				timeValue.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", apiResult.Qps), 64)
				qps := float64(resultData.Results[k].TotalRequestNum) / float64(resultData.Results[k].TotalRequestTime)
				resultData.Results[k].Qps, _ = strconv.ParseFloat(fmt.Sprintf("%0.2f", qps), 64)
				resultData.Results[k].QpsList = append(resultData.Results[k].QpsList, timeValue)
				timeValue.Value = resultData.Results[k].ErrorNum
				resultData.Results[k].ErrorNumList = append(resultData.Results[k].ErrorNumList, timeValue)
			}
		}
	}
	return

}

type SceneTestResultDataMsg struct {
	End        bool                             `json:"end" bson:"end"`
	ReportId   string                           `json:"report_id" bson:"report_id"`
	ReportName string                           `json:"report_name" bson:"report_name"`
	PlanId     int64                            `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName   string                           `json:"plan_name" bson:"plan_name"` //
	SceneId    int64                            `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName  string                           `json:"scene_name" bson:"scene_name"`
	Results    map[string]*ApiTestResultDataMsg `json:"results" bson:"results"`
	Machine    map[string]int64                 `json:"machine" bson:"machine"`
	TimeStamp  int64                            `json:"time_stamp" bson:"time_stamp"`
}

// ApiTestResultDataMsg 接口测试数据经过计算后的测试结果
type ApiTestResultDataMsg struct {
	Name                       string  `json:"name" bson:"name"`
	Concurrency                int64   `json:"concurrency" bson:"concurrency"`
	TotalRequestNum            uint64  `json:"total_request_num" bson:"total_request_num"`   // 总请求数
	TotalRequestTime           uint64  `json:"total_request_time" bson:"total_request_time"` // 总响应时间
	SuccessNum                 uint64  `json:"success_num" bson:"success_num"`
	ErrorNum                   uint64  `json:"error_num" bson:"error_num"`               // 错误数
	AvgRequestTime             uint64  `json:"avg_request_time" bson:"avg_request_time"` // 平均响应事件
	MaxRequestTime             uint64  `json:"max_request_time" bson:"max_request_time"`
	MinRequestTime             uint64  `json:"min_request_time" bson:"min_request_time"` // 毫秒
	CustomRequestTimeLine      uint64  `json:"custom_request_time_line" bson:"custom_request_time_line"`
	CustomRequestTimeLineValue int64   `json:"custom_request_time_line_value" bson:"custom_request_time_line_value"`
	NinetyRequestTimeLine      uint64  `json:"ninety_request_time_line" bson:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine  uint64  `json:"ninety_five_request_time_line" bson:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine  uint64  `json:"ninety_nine_request_time_line" bson:"ninety_nine_request_time_line"`
	SendBytes                  uint64  `json:"send_bytes" bson:"send_bytes"`         // 发送字节数
	ReceivedBytes              uint64  `json:"received_bytes" bson:"received_bytes"` // 接收字节数
	Qps                        float64 `json:"qps" bson:"qps"`
}

// ResultDataMsg 前端展示各个api数据
type ResultDataMsg struct {
	ApiName                    string      `json:"api_name" bson:"api_name"`
	Concurrency                int64       `json:"concurrency" bson:"concurrency"`
	TotalRequestNum            uint64      `json:"total_request_num" bson:"total_request_num"`   // 总请求数
	TotalRequestTime           uint64      `json:"total_request_time" bson:"total_request_time"` // 总响应时间
	SuccessNum                 uint64      `json:"success_num" bson:"success_num"`
	ErrorRate                  float64     `json:"error_rate" bson:"error_rate"`
	ErrorNum                   uint64      `json:"error_num" bson:"error_num"`               // 错误数
	AvgRequestTime             uint64      `json:"avg_request_time" bson:"avg_request_time"` // 平均响应事件
	MaxRequestTime             uint64      `json:"max_request_time" bson:"max_request_time"`
	MinRequestTime             uint64      `json:"min_request_time" bson:"min_request_time"` // 毫秒
	CustomRequestTimeLine      uint64      `json:"custom_request_time_line" bson:"custom_request_time_line"`
	CustomRequestTimeLineValue int64       `json:"custom_request_time_line_value" bson:"custom_request_time_line_value"`
	NinetyRequestTimeLine      uint64      `json:"ninety_request_time_line" bson:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine  uint64      `json:"ninety_five_request_time_line" bson:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine  uint64      `json:"ninety_nine_request_time_line" bson:"ninety_nine_request_time_line"`
	SendBytes                  uint64      `json:"send_bytes" bson:"send_bytes"`         // 发送字节数
	ReceivedBytes              uint64      `json:"received_bytes" bson:"received_bytes"` // 接收字节数
	Qps                        float64     `json:"qps" bson:"qps"`
	QpsList                    []TimeValue `json:"qps_list" bson:"qps_list"`
	ErrorNumList               []TimeValue `json:"error_rate_list" bson:"error_rate_list"`
}

type ResultData struct {
	End             bool                      `json:"end" bson:"end"`
	ReportId        string                    `json:"report_id" bson:"report_id"`
	ReportName      string                    `json:"report_name" bson:"report_name"`
	PlanId          int64                     `json:"plan_id" bson:"plan_id"`     // 任务ID
	PlanName        string                    `json:"plan_name" bson:"plan_name"` //
	SceneId         int64                     `json:"scene_id" bson:"scene_id"`   // 场景
	SceneName       string                    `json:"scene_name" bson:"scene_name"`
	Results         map[string]*ResultDataMsg `json:"results" bson:"results"`
	Machine         map[string]int64          `json:"machine" bson:"machine"`
	ConcurrencyList []TimeValue               `json:"concurrency_list" bson:"concurrency_list"`
	TimeStamp       int64                     `json:"time_stamp" bson:"time_stamp"`
}

type TimeValue struct {
	TimeStamp int64       `json:"time_stamp" bson:"time_stamp"`
	Value     interface{} `json:"value" bson:"value"`
}
