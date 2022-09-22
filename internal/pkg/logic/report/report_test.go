package report

import (
	"context"
	"encoding/json"
	"fmt"
	"kp-management/internal/pkg/dal/rao"
	"log"
	"testing"
)

func TestGetReportDetail(t *testing.T) {
	//client, _ := elastic.NewClient(
	//	elastic.SetURL("http://172.17.101.191:9200"),
	//	elastic.SetSniff(false),
	//	elastic.SetBasicAuth("elastic", "ZSrfx4R6ICa3skGBpCdf"),
	//	elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
	//	elastic.SetHealthcheckInterval(30*time.Second),
	//)
	//_, _, err := client.Ping("http://172.17.101.191:9200").Do(context.Background())
	//if err != nil {
	//	panic(fmt.Sprintf("es连接失败: %s", err))
	//}
	var report rao.GetReportReq
	report.ReportID = 762
	err, result := GetReportDetail(context.Background(), report, "http://172.17.101.191:9200", "elastic", "ZSrfx4R6ICa3skGBpCdf")
	if err != nil {
		fmt.Println(err)
	}
	res, _ := json.Marshal(result)
	log.Println(string(res))

}
