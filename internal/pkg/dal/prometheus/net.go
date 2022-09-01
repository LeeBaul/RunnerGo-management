package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"

	"kp-management/internal/pkg/conf"
)

func GetNetIORangeUsage(ip string, s, e int64) ([][]interface{}, error) {
	u := url.URL{
		Scheme:   "http",
		Host:     fmt.Sprintf("%s:%d", conf.Conf.Prometheus.Host, conf.Conf.Prometheus.Port),
		Path:     "/api/v1/query_range",
		RawQuery: "start=1661911200&end=1661911500&step=15&query=avg(rate(node_network_transmit_bytes_total{instance=\"172.17.101.188:9100\"}[1m]))",
	}

	uu := u.String()
	statusCode, body, err := fasthttp.Get(nil, uu)
	if err != nil {

	}
	if statusCode != http.StatusOK {

	}

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {

	}

	return resp.Data.Result[0].Values, nil
}
