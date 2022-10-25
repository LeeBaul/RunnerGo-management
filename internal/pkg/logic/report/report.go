package report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"time"

	"gorm.io/gen/field"

	"github.com/go-omnibus/proof"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gen"
	"kp-management/internal/pkg/biz/record"

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

func ListByTeamID2(ctx context.Context, teamID int64, limit, offset int, keyword string, startTimeSec, endTimeSec int64, taskType, taskMode, status, sortTag int32) ([]*rao.Report, int64, error) {

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

	if taskType > 0 {
		conditions = append(conditions, tx.TaskType.Eq(taskType))
	}

	if taskMode > 0 {
		conditions = append(conditions, tx.TaskMode.Eq(taskMode))
	}

	if status > 0 {
		conditions = append(conditions, tx.Status.Eq(status))
	}

	sort := make([]field.Expr, 0)

	if sortTag == 0 { // 默认排序
		sort = append(sort, tx.Rank.Desc())
		sort = append(sort, tx.ID.Desc())
	}
	if sortTag == 1 { // 创建时间倒序
		sort = append(sort, tx.CreatedAt.Desc())
	}
	if sortTag == 2 { // 创建时间正序
		sort = append(sort, tx.CreatedAt)
	}
	if sortTag == 3 { // 修改时间倒序
		sort = append(sort, tx.UpdatedAt.Desc())
	}
	if sortTag == 4 { // 修改时间正序
		sort = append(sort, tx.UpdatedAt)
	}

	reports, cnt, err := tx.WithContext(ctx).Where(conditions...).
		Order(sort...).
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

	return packer.TransReportModelToRaoReportList(reports, users), cnt, nil
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

func DeleteReport(ctx context.Context, teamID, reportID, userID int64) error {
	//tx := query.Use(dal.DB()).Report
	//_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(reportID)).Delete()
	//
	//return err

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		r, err := tx.Report.WithContext(ctx).Where(tx.Report.ID.Eq(reportID)).First()
		if err != nil {
			return err
		}

		if _, err := tx.Report.WithContext(ctx).Where(tx.Report.TeamID.Eq(teamID), tx.Report.ID.Eq(reportID)).Delete(); err != nil {
			return err
		}

		return record.InsertDelete(ctx, teamID, userID, record.OperationOperateDeleteReport, fmt.Sprintf("%s %s", r.PlanName, r.SceneName))
	})
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
		SceneID:        ru.SceneID,
		SceneName:      ru.SceneName,
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

func GetReportDebugStatus(ctx context.Context, report rao.GetReportReq) string {
	reportId := int(report.ReportID)
	//reportId := report.ReportID
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectDebugStatus)
	filter := bson.D{{"report_id", reportId}}
	//fmt.Println("filter:", filter)
	cur := collection.FindOne(ctx, filter)
	result, err := cur.DecodeBytes()
	if err != nil {
		return consts.StopDebug
	}
	list, err := result.Elements()
	if err != nil {
		return consts.StopDebug
	}
	for _, value := range list {
		if value.Key() == "debug" {
			return value.Value().StringValue()
		}
	}
	return consts.StopDebug
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

// GetReportDetail 从redis获取测试数据
func GetReportDetail(ctx context.Context, report rao.GetReportReq) (err error, resultData ResultData) {
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectReportData)
	filter := bson.D{{"reportid", fmt.Sprintf("%d", report.ReportID)}}
	var resultMsg SceneTestResultDataMsg

	var dataMap = make(map[string]string)
	err = collection.FindOne(ctx, filter).Decode(dataMap)
	if err != nil {
		rdb := dal.GetRDB()
		key := fmt.Sprintf("%d:%d:reportData", report.PlanId, report.ReportID)
		dataList := rdb.LRange(ctx, key, 0, -1).Val()
		if len(dataList) < 0 {
			return
		}
		fmt.Println(len(dataList))
		for i := len(dataList) - 1; i >= 0; i-- {
			resultMsgString := dataList[i]
			err = json.Unmarshal([]byte(resultMsgString), &resultMsg)
			if err != nil {
				proof.Error("json转换格式错误：", proof.WithError(err))
			}
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
					resultData.Results[k].TotalRequestTime, _ = decimal.NewFromFloat(float64(apiResult.TotalRequestTime) / float64(time.Second)).Round(2).Float64()
					resultData.Results[k].SuccessNum = apiResult.SuccessNum
					resultData.Results[k].ErrorNum = apiResult.ErrorNum
					if resultData.Results[k].ErrorNum != 0 && apiResult.TotalRequestNum != 0 {
						errRate := float64(apiResult.ErrorNum) / float64(apiResult.TotalRequestNum)
						resultData.Results[k].ErrorRate, _ = decimal.NewFromFloat(errRate).Round(4).Float64()
					}

					resultData.Results[k].AvgRequestTime, _ = decimal.NewFromFloat(apiResult.AvgRequestTime / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].MaxRequestTime, _ = decimal.NewFromFloat(apiResult.MaxRequestTime / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].MinRequestTime, _ = decimal.NewFromFloat(apiResult.MinRequestTime / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].CustomRequestTimeLine = apiResult.CustomRequestTimeLine
					resultData.Results[k].CustomRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.CustomRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].NinetyRequestTimeLine = apiResult.NinetyRequestTimeLine
					resultData.Results[k].NinetyRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.NinetyRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].NinetyFiveRequestTimeLine = apiResult.NinetyFiveRequestTimeLine
					resultData.Results[k].NinetyFiveRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.NinetyFiveRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].NinetyNineRequestTimeLine = apiResult.NinetyNineRequestTimeLine
					resultData.Results[k].NinetyNineRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.NinetyNineRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
					resultData.Results[k].SendBytes, _ = decimal.NewFromFloat(apiResult.SendBytes).Round(1).Float64()
					resultData.Results[k].ReceivedBytes, _ = decimal.NewFromFloat(apiResult.ReceivedBytes).Round(1).Float64()
					resultData.Results[k].Qps = apiResult.Qps
					if resultData.Results[k].QpsList == nil {
						resultData.Results[k].QpsList = []TimeValue{}
					}
					var timeValue = TimeValue{}
					timeValue.TimeStamp = resultData.TimeStamp
					// qps列表
					timeValue.Value = resultData.Results[k].Qps
					resultData.Results[k].QpsList = append(resultData.Results[k].QpsList, timeValue)
					timeValue.Value = resultData.Results[k].ErrorNum
					if resultData.Results[k].ErrorNumList == nil {
						resultData.Results[k].ErrorNumList = []TimeValue{}
					}
					// 错误数列表
					resultData.Results[k].ErrorNumList = append(resultData.Results[k].ErrorNumList, timeValue)
					timeValue.Value = resultData.Results[k].Concurrency
					if resultData.Results[k].ConcurrencyList == nil {
						resultData.Results[k].ConcurrencyList = []TimeValue{}
					}
					// 并发数列表
					resultData.Results[k].ConcurrencyList = append(resultData.Results[k].ConcurrencyList, timeValue)

					// 平均响应时间列表
					timeValue.Value = resultData.Results[k].AvgRequestTime
					if resultData.Results[k].AvgList == nil {
						resultData.Results[k].AvgList = []TimeValue{}
					}
					resultData.Results[k].AvgList = append(resultData.Results[k].AvgList, timeValue)

					// 90响应时间列表
					timeValue.Value = resultData.Results[k].NinetyNineRequestTimeLineValue
					if resultData.Results[k].NinetyList == nil {
						resultData.Results[k].NinetyList = []TimeValue{}
					}
					resultData.Results[k].NinetyList = append(resultData.Results[k].NinetyList, timeValue)

					// 95响应时间列表
					timeValue.Value = resultData.Results[k].NinetyFiveRequestTimeLineValue
					if resultData.Results[k].NinetyFiveList == nil {
						resultData.Results[k].NinetyFiveList = []TimeValue{}
					}
					resultData.Results[k].NinetyFiveList = append(resultData.Results[k].NinetyFiveList, timeValue)

					// 99响应时间列表
					timeValue.Value = resultData.Results[k].NinetyNineRequestTimeLineValue
					if resultData.Results[k].NinetyNineList == nil {
						resultData.Results[k].NinetyNineList = []TimeValue{}
					}
					resultData.Results[k].NinetyNineList = append(resultData.Results[k].NinetyNineList, timeValue)
				}
			}
			if resultMsg.End {
				var by []byte
				by, err = json.Marshal(resultData)
				if err != nil {
					proof.Error("resultData转字节失败：：    ", proof.WithError(err))
					return
				}
				dataMap["reportid"] = resultData.ReportId
				dataMap["data"] = string(by)
				_, err = collection.InsertOne(ctx, dataMap)
				if err != nil {
					proof.Error("测试数据写入mongo失败：    ", proof.WithError(err))
					return
				}
				err = rdb.Del(ctx, key).Err()
				if err != nil {
					proof.Error(fmt.Sprintf("删除redis的key：%s:    ", key), proof.WithError(err))
					return
				}
			}
		}
	} else {
		data := dataMap["data"]
		json.Unmarshal([]byte(data), &resultData)
		return
	}
	return
}

// 从es获取测试数据
//func GetReportDetail(ctx context.Context, report rao.GetReportReq, host, user, password string) (err error, resultData ResultData) {
//	reportId := strconv.FormatInt(report.ReportID, 10)
//	//index := strconv.FormatInt(report.TeamID, 10)
//
//	queryEs := elastic.NewBoolQuery()
//	queryEs = queryEs.Must(elastic.NewMatchQuery("report_id", reportId))
//
//	client, err := elastic.NewClient(
//		elastic.SetURL(host),
//		elastic.SetSniff(false),
//		elastic.SetBasicAuth(user, password),
//		elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
//		elastic.SetHealthcheckInterval(30*time.Second),
//	)
//	if err != nil {
//		if err != nil {
//			proof.Error("创建es客户端失败", proof.WithError(err))
//			return
//		}
//	}
//	_, _, err = client.Ping(host).Do(ctx)
//	if err != nil {
//		proof.Error("es连接失败", proof.WithError(err))
//		return
//	}
//	res, err := client.Search(reportId).Query(queryEs).Sort("time_stamp", true).Size(conf.Conf.ES.Size).Pretty(true).Do(ctx)
//	if err != nil {
//		proof.Error("获取报告详情失败", proof.WithError(err))
//		return
//	}
//	if res == nil {
//		proof.Error("报告详情为空")
//		return
//	}
//
//	var resultMsg SceneTestResultDataMsg // 从es中获取得数据结构
//
//	for _, item := range res.Hits.Hits {
//		err = json.Unmarshal(item.Source, &resultMsg)
//		if err != nil {
//			proof.Error("json转换格式错误：", proof.WithError(err))
//		}
//		if resultData.Results == nil {
//			resultData.Results = make(map[string]*ResultDataMsg)
//		}
//		resultData.ReportId = resultMsg.ReportId
//		resultData.End = resultMsg.End
//		resultData.ReportName = resultMsg.ReportName
//		resultData.PlanId = resultMsg.PlanId
//		resultData.PlanName = resultMsg.PlanName
//		resultData.SceneId = resultMsg.SceneId
//		resultData.SceneName = resultMsg.SceneName
//		resultData.TimeStamp = resultMsg.TimeStamp
//		if resultMsg.Results != nil && len(resultMsg.Results) > 0 {
//			for k, apiResult := range resultMsg.Results {
//				if resultData.Results[k] == nil {
//					resultData.Results[k] = new(ResultDataMsg)
//				}
//				resultData.Results[k].ApiName = apiResult.Name
//				resultData.Results[k].Concurrency = apiResult.Concurrency
//				resultData.Results[k].TotalRequestNum = apiResult.TotalRequestNum
//				resultData.Results[k].TotalRequestTime, _ = decimal.NewFromFloat(float64(apiResult.TotalRequestTime) / float64(time.Second)).Round(2).Float64()
//				resultData.Results[k].SuccessNum = apiResult.SuccessNum
//				resultData.Results[k].ErrorNum = apiResult.ErrorNum
//				if resultData.Results[k].ErrorNum != 0 && apiResult.TotalRequestNum != 0 {
//					errRate := float64(apiResult.ErrorNum) / float64(apiResult.TotalRequestNum)
//					resultData.Results[k].ErrorRate, _ = decimal.NewFromFloat(errRate).Round(4).Float64()
//				}
//
//				resultData.Results[k].AvgRequestTime, _ = decimal.NewFromFloat(apiResult.AvgRequestTime / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].MaxRequestTime, _ = decimal.NewFromFloat(apiResult.MaxRequestTime / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].MinRequestTime, _ = decimal.NewFromFloat(apiResult.MinRequestTime / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].CustomRequestTimeLine = apiResult.CustomRequestTimeLine
//				resultData.Results[k].CustomRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.CustomRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].NinetyRequestTimeLine = apiResult.NinetyRequestTimeLine
//				resultData.Results[k].NinetyRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.NinetyRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].NinetyFiveRequestTimeLine = apiResult.NinetyFiveRequestTimeLine
//				resultData.Results[k].NinetyFiveRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.NinetyFiveRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].NinetyNineRequestTimeLine = apiResult.NinetyNineRequestTimeLine
//				resultData.Results[k].NinetyNineRequestTimeLineValue, _ = decimal.NewFromFloat(apiResult.NinetyNineRequestTimeLineValue / float64(time.Millisecond)).Round(1).Float64()
//				resultData.Results[k].SendBytes, _ = decimal.NewFromFloat(apiResult.SendBytes).Round(1).Float64()
//				resultData.Results[k].ReceivedBytes, _ = decimal.NewFromFloat(apiResult.ReceivedBytes).Round(1).Float64()
//				resultData.Results[k].Qps = apiResult.Qps
//				if resultData.Results[k].QpsList == nil {
//					resultData.Results[k].QpsList = []TimeValue{}
//				}
//				var timeValue = TimeValue{}
//				timeValue.TimeStamp = resultData.TimeStamp
//				// qps列表
//				timeValue.Value = resultData.Results[k].Qps
//				resultData.Results[k].QpsList = append(resultData.Results[k].QpsList, timeValue)
//				timeValue.Value = resultData.Results[k].ErrorNum
//				if resultData.Results[k].ErrorNumList == nil {
//					resultData.Results[k].ErrorNumList = []TimeValue{}
//				}
//				// 错误数列表
//				resultData.Results[k].ErrorNumList = append(resultData.Results[k].ErrorNumList, timeValue)
//				timeValue.Value = resultData.Results[k].Concurrency
//				if resultData.Results[k].ConcurrencyList == nil {
//					resultData.Results[k].ConcurrencyList = []TimeValue{}
//				}
//				// 并发数列表
//				resultData.Results[k].ConcurrencyList = append(resultData.Results[k].ConcurrencyList, timeValue)
//
//				// 平均响应时间列表
//				timeValue.Value = resultData.Results[k].AvgRequestTime
//				if resultData.Results[k].AvgList == nil {
//					resultData.Results[k].AvgList = []TimeValue{}
//				}
//				resultData.Results[k].AvgList = append(resultData.Results[k].AvgList, timeValue)
//
//				// 90响应时间列表
//				timeValue.Value = resultData.Results[k].NinetyNineRequestTimeLineValue
//				if resultData.Results[k].NinetyList == nil {
//					resultData.Results[k].NinetyList = []TimeValue{}
//				}
//				resultData.Results[k].NinetyList = append(resultData.Results[k].NinetyList, timeValue)
//
//				// 95响应时间列表
//				timeValue.Value = resultData.Results[k].NinetyFiveRequestTimeLineValue
//				if resultData.Results[k].NinetyFiveList == nil {
//					resultData.Results[k].NinetyFiveList = []TimeValue{}
//				}
//				resultData.Results[k].NinetyFiveList = append(resultData.Results[k].NinetyFiveList, timeValue)
//
//				// 99响应时间列表
//				timeValue.Value = resultData.Results[k].NinetyNineRequestTimeLineValue
//				if resultData.Results[k].NinetyNineList == nil {
//					resultData.Results[k].NinetyNineList = []TimeValue{}
//				}
//				resultData.Results[k].NinetyNineList = append(resultData.Results[k].NinetyNineList, timeValue)
//
//			}
//		}
//	}
//	return
//
//}

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
	Name                           string  `json:"name" bson:"name"`
	Concurrency                    int64   `json:"concurrency" bson:"concurrency"`
	TotalRequestNum                uint64  `json:"total_request_num" bson:"total_request_num"`   // 总请求数
	TotalRequestTime               uint64  `json:"total_request_time" bson:"total_request_time"` // 总响应时间
	SuccessNum                     uint64  `json:"success_num" bson:"success_num"`
	ErrorNum                       uint64  `json:"error_num" bson:"error_num"`               // 错误数
	ErrorThreshold                 float64 `json:"error_threshold" bson:"error_threshold"`   // 自定义错误率
	AvgRequestTime                 float64 `json:"avg_request_time" bson:"avg_request_time"` // 平均响应事件
	MaxRequestTime                 float64 `json:"max_request_time" bson:"max_request_time"`
	MinRequestTime                 float64 `json:"min_request_time" bson:"min_request_time"` // 毫秒
	CustomRequestTimeLine          int64   `json:"custom_request_time_line" bson:"custom_request_time_line"`
	NinetyRequestTimeLine          int64   `json:"ninety_request_time_line" bson:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine      int64   `json:"ninety_five_request_time_line" bson:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine      int64   `json:"ninety_nine_request_time_line" bson:"ninety_nine_request_time_line"`
	CustomRequestTimeLineValue     float64 `json:"custom_request_time_line_value" bson:"custom_request_time_line_value"`
	NinetyRequestTimeLineValue     float64 `json:"ninety_request_time_line_value" bson:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64 `json:"ninety_five_request_time_line_value" bson:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64 `json:"ninety_nine_request_time_line_value" bson:"ninety_nine_request_time_line_value"`
	SendBytes                      float64 `json:"send_bytes" bson:"send_bytes"`         // 发送字节数
	ReceivedBytes                  float64 `json:"received_bytes" bson:"received_bytes"` // 接收字节数
	Qps                            float64 `json:"qps" bson:"qps"`
}

// ResultDataMsg 前端展示各个api数据
type ResultDataMsg struct {
	ApiName                        string      `json:"api_name" bson:"api_name"`
	Concurrency                    int64       `json:"concurrency" bson:"concurrency"`
	TotalRequestNum                uint64      `json:"total_request_num" bson:"total_request_num"`   // 总请求数
	TotalRequestTime               float64     `json:"total_request_time" bson:"total_request_time"` // 总响应时间
	SuccessNum                     uint64      `json:"success_num" bson:"success_num"`
	ErrorRate                      float64     `json:"error_rate" bson:"error_rate"`
	ErrorNum                       uint64      `json:"error_num" bson:"error_num"`               // 错误数
	AvgRequestTime                 float64     `json:"avg_request_time" bson:"avg_request_time"` // 平均响应事件
	MaxRequestTime                 float64     `json:"max_request_time" bson:"max_request_time"`
	MinRequestTime                 float64     `json:"min_request_time" bson:"min_request_time"` // 毫秒
	CustomRequestTimeLine          int64       `json:"custom_request_time_line" bson:"custom_request_time_line"`
	NinetyRequestTimeLine          int64       `json:"ninety_request_time_line" bson:"ninety_request_time_line"`
	NinetyFiveRequestTimeLine      int64       `json:"ninety_five_request_time_line" bson:"ninety_five_request_time_line"`
	NinetyNineRequestTimeLine      int64       `json:"ninety_nine_request_time_line" bson:"ninety_nine_request_time_line"`
	CustomRequestTimeLineValue     float64     `json:"custom_request_time_line_value" bson:"custom_request_time_line_value"`
	NinetyRequestTimeLineValue     float64     `json:"ninety_request_time_line_value" bson:"ninety_request_time_line_value"`
	NinetyFiveRequestTimeLineValue float64     `json:"ninety_five_request_time_line_value" bson:"ninety_five_request_time_line_value"`
	NinetyNineRequestTimeLineValue float64     `json:"ninety_nine_request_time_line_value" bson:"ninety_nine_request_time_line_value"`
	SendBytes                      float64     `json:"send_bytes" bson:"send_bytes"`         // 发送字节数
	ReceivedBytes                  float64     `json:"received_bytes" bson:"received_bytes"` // 接收字节数
	Qps                            float64     `json:"qps" bson:"qps"`
	ConcurrencyList                []TimeValue `json:"concurrency_list"`
	QpsList                        []TimeValue `json:"qps_list" bson:"qps_list"`
	ErrorNumList                   []TimeValue `json:"error_num_list" bson:"error_num_list"`
	AvgList                        []TimeValue `json:"avg_list" bson:"avg_list"`
	NinetyList                     []TimeValue `json:"ninety_list" bson:"ninety_list"`
	NinetyFiveList                 []TimeValue `json:"ninety_five_list" bson:"ninety_five_list"`
	NinetyNineList                 []TimeValue `json:"ninety_nine_list" bson:"ninety_nine_list"`
}

type ResultData struct {
	End          bool                      `json:"end"`
	ReportId     string                    `json:"report_id"`
	ReportName   string                    `json:"report_name"`
	PlanId       int64                     `json:"plan_id"`   // 任务ID
	PlanName     string                    `json:"plan_name"` //
	SceneId      int64                     `json:"scene_id"`  // 场景
	SceneName    string                    `json:"scene_name"`
	Results      map[string]*ResultDataMsg `json:"results"`
	TotalQps     []float64                 `json:"total_qps"`
	Machine      map[string]int64          `json:"machine"`
	TotalQpsList []TimeValue               `json:"total_qps_list"`
	TimeStamp    int64                     `json:"time_stamp"`
}

type TimeValue struct {
	TimeStamp int64       `json:"time_stamp" bson:"time_stamp"`
	Value     interface{} `json:"value" bson:"value"`
}
