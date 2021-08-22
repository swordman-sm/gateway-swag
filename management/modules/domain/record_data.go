package domain

type RecordsData struct {
	Time        int64                             `json:"time"`
	MetricsData map[string]map[string]interface{} `json:"metrics_data"`
}
