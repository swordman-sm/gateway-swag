package modules

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"net/url"
	"strconv"
	"time"
)

var LbMap = map[string]bool{"roundRobin": true, "random": true}

type Target struct {
	Pointer       string `json:"pointer"`
	Weight        int8   `json:"weight"`
	CurrentWeight int8   `json:"current_weight"`
}

type Domain struct {
	Id         string `json:"id"`
	DomainName string `json:"domain_name"`
	DomainUrl  string `json:"domain_url"`
	LbType     string `json:"lb_type"`
	//代理实现
	Targets            []*Target       `json:"targets"`
	BlackIps           map[string]bool `json:"black_ips"`
	RateLimiterNum     float64         `json:"rate_limiter_num"`
	RateLimiterMsg     string          `json:"rate_limiter_msg"`
	RateLimiterEnabled bool            `json:"rate_limiter_enabled"`
	SetTime            string          `json:"set_time"`
}

func QueryAllDomains(ctx *gin.Context) {
	resp, err := getAllDomainsData()
	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}

	var domains []*Domain
	if resp.Count > 0 {
		for _, kv := range resp.Kvs {
			domain := new(Domain)
			err := json.Unmarshal(kv.Value, domain)
			if err == nil {
				domains = append(domains, domain)
			}
		}
		resultCtx{ctx}.SucResult(domains)
		return
	}
	resultCtx{ctx}.SucResult(make([]string, 0))
}

func QueryDomainDataByDomainId(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	if domainId == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}
	rsp, err := getDomainDataByDomainId(domainId)
	if err != nil {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	if rsp.Count > 0 {
		domain := new(Domain)
		err := json.Unmarshal(rsp.Kvs[0].Value, domain)
		if err != nil {
			resultCtx{ctx}.ErrResult(DataParseError)
			return
		}
		resultCtx{ctx}.SucResult(domain)
	} else {
		resultCtx{ctx}.SucResult(struct{}{})
	}
}

func DelDomainByDomainId(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	if domainId == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	deleted := delDomainDataByDomainId(domainId)
	if deleted {
		resultCtx{ctx}.SucResult(nil)
		return
	}
	resultCtx{ctx}.ErrResult(SystemError)
}

func AddDomainData(ctx *gin.Context) {
	domainUrl := ctx.PostForm("domain_url")
	domainName := ctx.PostForm("domain_name")
	lbType := ctx.PostForm("lb_type")
	proxyTargets := ctx.PostForm("proxy_targets")
	blackIpsJson := ctx.PostForm("black_ips")
	rateLimiterNum := ctx.PostForm("rate_limiter_num")
	rateLimiterMsg := ctx.PostForm("rate_limiter_msg")
	rateLimiterEnabled := ctx.PostForm("rate_limiter_enabled")

	if domainUrl == "" || domainName == "" || lbType == "" || proxyTargets == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	urlParse, err := url.ParseRequestURI(domainUrl)
	if err != nil {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	//检查负载均衡模式
	if _, ok := LbMap[lbType]; !ok {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	//检查代理目标数据
	var targets []*Target
	err = json.Unmarshal([]byte(proxyTargets), &targets)
	if err != nil {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	//黑名单解析
	blackIps := make(map[string]bool)
	if blackIpsJson != "" {
		err := json.Unmarshal([]byte(blackIpsJson), &blackIps)
		if err != nil {
			resultCtx{ctx}.ErrResult(DataParseError)
			return
		}
	}

	//有接收到domainId 就是修改操作， 否则就是新增
	var domainId string
	domainId = ctx.Param("domain_id")
	if domainId == "" {
		domainId = uuid.Must(uuid.NewV4()).String()
	}
	domain := new(Domain)
	domain.Id = domainId
	domain.DomainName = domainName
	domain.DomainUrl = urlParse.Host
	domain.LbType = lbType
	domain.Targets = targets
	domain.BlackIps = blackIps
	domain.RateLimiterNum, _ = strconv.ParseFloat(rateLimiterNum, 10)
	domain.RateLimiterMsg = rateLimiterMsg
	domain.RateLimiterEnabled, _ = strconv.ParseBool(rateLimiterEnabled)
	domain.SetTime = time.Now().Format("2006/1/2 15:04:05")

	domainB, err := json.Marshal(domain)
	if err != nil {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	err = addDomainData(domain.Id, string(domainB))
	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	resultCtx{ctx}.SucResult(domain)
}
