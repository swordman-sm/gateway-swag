package modules

import (
	"encoding/json"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
)

func AddRequestListen(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	listenPath := ctx.PostForm("listen_path")
	if domainId == "" || listenPath == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	domainRsp, err := getDomainDataByDomainId(domainId)
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
	err = putRequestListen(listenId, string(reqCopyB))
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	base.Result{Context: ctx}.SucResult(make([]string, 0))
}

func RequestsCopy(ctx *gin.Context) {
	var data []*domain.RequestCopy
	rsp, err := requestsCopy()
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
