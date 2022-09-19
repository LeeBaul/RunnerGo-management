package report

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"testing"
	"time"
)

func TestGetReportDetail(t *testing.T) {
	client, _ := elastic.NewClient(
		elastic.SetURL("http://172.17.101.191:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "ZSrfx4R6ICa3skGBpCdf"),
		elastic.SetErrorLog(log.New(os.Stdout, "APP", log.Lshortfile)),
		elastic.SetHealthcheckInterval(30*time.Second),
	)
	_, _, err := client.Ping("http://172.17.101.191:9200").Do(context.Background())
	if err != nil {
		panic(fmt.Sprintf("es连接失败: %s", err))
	}

	////删除指定index中的所有数据
	//_, err = client.DeleteIndex("report").Do(context.Background())
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//result := GetReportDetail("report", "2222222",)
	//res, _ := json.Marshal(result)
	//log.Println(string(res))

}
