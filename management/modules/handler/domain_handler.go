package handler

import (
	"encoding/json"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"gateway-swag/management/modules/service/impl"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"net/url"
	"strconv"
	"time"
)

var LbMap = map[string]bool{"roundRobin": true, "random": true}

var domainService = new(impl.DomainServiceImpl)

func QueryAllDomainsDataHandler(ctx *gin.Context) {
	resp, err := domainService.GetAllDomainsData()
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}

	var domains []*domain.Domain
	if resp.Count > 0 {
		for _, kv := range resp.Kvs {
			domain := new(domain.Domain)
			err := json.Unmarshal(kv.Value, domain)
			if err == nil {
				domains = append(domains, domain)
			}
		}
		base.Result{Context: ctx}.SucResult(domains)
		return
	}
	base.Result{Context: ctx}.SucResult(make([]string, 0))
}

func QueryDomainDataByDomainIdHandler(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	if domainId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}
	rsp, err := domainService.GetDomainDataByDomainId(domainId)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	if rsp.Count > 0 {
		domain := new(domain.Domain)
		err := json.Unmarshal(rsp.Kvs[0].Value, domain)
		if err != nil {
			base.Result{Context: ctx}.ErrResult(base.DataParseError)
			return
		}
		base.Result{Context: ctx}.SucResult(domain)
	} else {
		base.Result{Context: ctx}.SucResult(struct{}{})
	}
}

func DelDomainByDomainIdHandler(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	if domainId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	deleted := domainService.DelDomainDataByDomainId(domainId)
	if deleted {
		base.Result{Context: ctx}.SucResult(nil)
		return
	}
	base.Result{Context: ctx}.ErrResult(base.SystemError)
}

func AddDomainDataHandler(ctx *gin.Context) {
	domainUrl := ctx.PostForm("domain_url")
	domainName := ctx.PostForm("domain_name")
	lbType := ctx.PostForm("lb_type")
	proxyTargets := ctx.PostForm("proxy_targets")
	blackIpsJson := ctx.PostForm("black_ips")
	rateLimiterNum := ctx.PostForm("rate_limiter_num")
	rateLimiterMsg := ctx.PostForm("rate_limiter_msg")
	rateLimiterEnabled := ctx.PostForm("rate_limiter_enabled")

	if domainUrl == "" || domainName == "" || lbType == "" || proxyTargets == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	urlParse, err := url.ParseRequestURI(domainUrl)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	//检查负载均衡模式
	if _, ok := LbMap[lbType]; !ok {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	//检查代理目标数据
	var targets []*domain.Target
	err = json.Unmarshal([]byte(proxyTargets), &targets)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	//黑名单解析
	blackIps := make(map[string]bool)
	if blackIpsJson != "" {
		err := json.Unmarshal([]byte(blackIpsJson), &blackIps)
		if err != nil {
			base.Result{Context: ctx}.ErrResult(base.DataParseError)
			return
		}
	}

	//有接收到domainId 就是修改操作， 否则就是新增
	var domainId string
	domainId = ctx.Param("domain_id")
	if domainId == "" {
		domainId = uuid.Must(uuid.NewV4()).String()
	}
	domainA := new(domain.Domain)
	domainA.Id = domainId
	domainA.DomainName = domainName
	domainA.DomainUrl = urlParse.Host
	domainA.LbType = lbType
	domainA.Targets = targets
	domainA.BlackIps = blackIps
	domainA.RateLimiterNum, _ = strconv.ParseFloat(rateLimiterNum, 10)
	domainA.RateLimiterMsg = rateLimiterMsg
	domainA.RateLimiterEnabled, _ = strconv.ParseBool(rateLimiterEnabled)
	domainA.SetTime = time.Now().Format("2006/1/2 15:04:05")

	domainB, err := json.Marshal(domainA)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	err = domainService.AddDomainData(domainA.Id, string(domainB))
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	base.Result{Context: ctx}.SucResult(domainA)
}
