package domain

import "net/http"

type RequestCopy struct {
	SerName   string      `json:"ser_name"`
	Id        string      `json:"id"`
	ReqTime   string      `json:"req_time"`
	ReqIp     string      `json:"req_ip"`
	ReqPath   string      `json:"req_path"`
	PostForm  interface{} `json:"post_form"`
	Get       string      `json:"get"`
	ReqHeader interface{} `json:"req_header"`
	RspSize   int         `json:"rsp_size"`
	RspHeader http.Header `json:"rsp_header"`
	RspBody   string      `json:"rsp_body"`
}
