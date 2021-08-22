package handler

import (
	"encoding/json"
	"fmt"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"gateway-swag/management/modules/service/impl"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Gw struct {
	ServerName string `json:"server_name"`
}

type Record struct {
	Count   float64 `json:"count"`
	Mean    float64 `json:"mean"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	TimeStr string  `json:"time_str"`
}

const (
	hour          = 60
	halfHour      = 30
	fifteenMinute = 15
	fiveMinute    = 5
)

var gatewayService = new(impl.GatewayServiceImpl)

func GetAllGatewayDataHandler(ctx *gin.Context) {
	rsp, err := gatewayService.GetAllGatewaysData()
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	var gateways []*Gw
	if rsp.Count > 0 {
		for _, kv := range rsp.Kvs {
			gateway := new(Gw)
			gateway.ServerName = strings.Replace(string(kv.Key), base.MetricsGatewayActivePrefix, "", 1)
			gateways = append(gateways, gateway)
		}
		base.Result{Context: ctx}.SucResult(gateways)
		return
	}
	base.Result{Context: ctx}.SucResult(make([]string, 0))
}

func GetGatewayDataByServerHandler(ctx *gin.Context) {
	serName := ctx.Param("server_name")
	if serName == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}
	mcRsp, _ := gatewayService.GetGatewayDataByServer(serName)
	machineData := string(mcRsp.Kvs[0].Value)

	//取30次以内的结果
	rsp, err := gatewayService.GetGatewayDataByLimit(serName, 30)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	metricsMap := make(map[string][]Record)
	var timesData []string
	spanMap := make(map[string]string)

	for _, kv := range rsp.Kvs {
		recordData := new(domain.RecordsData)
		err := json.Unmarshal(kv.Value, &recordData)
		if err != nil {
			continue
		}
		fmt.Println(recordData)
		recordTime := recordData.Time
		timeData := time.Unix(recordTime, 0).Format("15:04")
		timesData = append(timesData, timeData)
		metricsData := recordData.MetricsData
		for oldSpan, _ := range metricsData {
			if _, ok := spanMap[oldSpan]; !ok {
				spans := strings.Split(oldSpan, "|-|")
				spanMap[oldSpan] = spans[0] + spans[1]
			}
			span := spanMap[oldSpan]
			if _, ok := metricsMap[span]; !ok {
				var data []Record
				metricsMap[span] = data
			}
		}
	}
	start := 0
	for _, kv := range rsp.Kvs {
		recordData := domain.RecordsData{}
		err := json.Unmarshal(kv.Value, &recordData)
		if err != nil {
			continue
		}
		recordTime := recordData.Time
		timeData := time.Unix(recordTime, 0).Format("15:04")
		for k := range metricsMap {
			metricsMap[k] = append(metricsMap[k], Record{0, 0, 0, 0, timeData})
		}
		for oldSpan, metrics := range recordData.MetricsData {
			span := spanMap[oldSpan]
			metricsMap[span][start].Count = metrics["count"].(float64)
			metricsMap[span][start].Mean = Milliseconds(time.Duration(int64(metrics["mean"].(float64))))
			metricsMap[span][start].Max = Milliseconds(time.Duration(int64(metrics["max"].(float64))))
			metricsMap[span][start].Min = Milliseconds(time.Duration(int64(metrics["min"].(float64))))
		}
		start++
	}
	allData := make(map[string]interface{})
	allData["metrics"] = metricsMap
	allData["times"] = timesData
	allData["machine"] = machineData
	base.Result{Context: ctx}.SucResult(allData)
}

func Milliseconds(d time.Duration) float64 {
	mill := d / time.Millisecond
	micrs := d % time.Microsecond
	return float64(mill) + float64(micrs)/1e6
}
