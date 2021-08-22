package handler

import (
	"encoding/json"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"gateway-swag/management/modules/service/impl"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
)

var requestListenService = new(impl.RequestListenServiceImpl)

func AddRequestListenHandler(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	listenPath := ctx.PostForm("listen_path")
	if domainId == "" || listenPath == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	domainRsp, err := domainService.GetDomainDataByDomainId(domainId)
	if err != nil || domainRsp.Count == 0 {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}
	dom := new(domain.Domain)
	err = json.Unmarshal(domainRsp.Kvs[0].Value, dom)
	if err != nil || dom == nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	reqListen := new(domain.RequestListen)
	reqListen.DomainUrl = dom.DomainUrl
	reqListen.ListenPath = listenPath
	reqCopyB, _ := json.Marshal(reqListen)
	listenId := uuid.Must(uuid.NewV4()).String()
	err = requestListenService.AddRequestListen(listenId, string(reqCopyB))
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	base.Result{Context: ctx}.SucResult(make([]string, 0))
}

func RequestsCopyHandler(ctx *gin.Context) {
	var data []*domain.RequestCopy
	rsp, err := requestListenService.GetRequestsCopy()
	if err != nil || rsp.Count == 0 {
		base.Result{Context: ctx}.SucResult(make([]string, 0))
		return
	}

	for _, kv := range rsp.Kvs {
		reqCopy := new(domain.RequestCopy)
		err := json.Unmarshal(kv.Value, reqCopy)
		if err != nil {
			continue
		}
		data = append(data, reqCopy)
	}
	base.Result{Context: ctx}.SucResult(data)
}
