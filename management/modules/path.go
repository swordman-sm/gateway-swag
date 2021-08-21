package modules

import (
	"encoding/json"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"strconv"
	"time"
)

func Paths(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	if domainId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	rsp, err := getPathDataByDomainId(domainId)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	var paths []*domain.Path
	if rsp.Count > 0 {
		for _, kv := range rsp.Kvs {
			path := new(domain.Path)
			err := json.Unmarshal(kv.Value, path)
			if err == nil {
				paths = append(paths, path)
			}
		}
		base.Result{Context: ctx}.SucResult(paths)
		return
	}
	base.Result{Context: ctx}.SucResult(make([]string, 0))
}

func PutPath(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	reqMethod := ctx.PostForm("req_method")
	reqPath := ctx.PostForm("req_path")
	searchPath := ctx.PostForm("search_path")
	replacePath := ctx.PostForm("replace_path")
	cbRequest := ctx.PostForm("circuit_breaker_request")
	cbPercent := ctx.PostForm("circuit_breaker_percent")
	cbTimeout := ctx.PostForm("circuit_breaker_timeout")
	cbMsg := ctx.PostForm("circuit_breaker_msg")
	cbEnabled := ctx.PostForm("circuit_breaker_enabled")
	cbForce := ctx.PostForm("circuit_breaker_force")
	priProxyEnabled := ctx.PostForm("private_proxy_enabled")
	lbType := ctx.PostForm("lb_type")
	proxyTargets := ctx.PostForm("proxy_targets")

	if reqMethod == "" || reqPath == "" || domainId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	//检查代理目标数据
	var targets []*domain.Target
	err := json.Unmarshal([]byte(proxyTargets), &targets)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	var pathId string
	pathId = ctx.Param("path_id")
	if pathId == "" {
		pathId = uuid.Must(uuid.NewV4()).String()
	}

	path := new(domain.Path)
	path.Id = pathId
	path.DomainId = domainId
	path.ReqMethod = reqMethod
	path.ReqPath = reqPath
	path.SearchPath = searchPath
	path.ReplacePath = replacePath
	path.CircuitBreakerRequest, _ = strconv.Atoi(cbRequest)
	path.CircuitBreakerPercent, _ = strconv.Atoi(cbPercent)
	path.CircuitBreakerTimeout, _ = strconv.Atoi(cbTimeout)
	path.CircuitBreakerMsg = cbMsg
	path.CircuitBreakerEnabled, _ = strconv.ParseBool(cbEnabled)
	path.CircuitBreakerForce, _ = strconv.ParseBool(cbForce)
	path.LbType = lbType
	path.Targets = targets
	path.PrivateProxyEnabled, _ = strconv.ParseBool(priProxyEnabled)
	path.SetTime = time.Now().Format("2006/1/2 15:04:05")

	pathB, err := json.Marshal(path)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	err = addPathData(domainId, path.Id, string(pathB))
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	base.Result{Context: ctx}.SucResult(path)
}

func GetPath(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	pathId := ctx.Param("path_id")
	if domainId == "" || pathId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}
	rsp, err := getPathDataByDomainIdAndPathId(domainId, pathId)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	if rsp.Count > 0 {
		path := new(domain.Path)
		err := json.Unmarshal(rsp.Kvs[0].Value, path)
		if err != nil {
			base.Result{Context: ctx}.ErrResult(base.DataParseError)
			return
		}
		base.Result{Context: ctx}.SucResult(path)
	} else {
		base.Result{Context: ctx}.SucResult(struct{}{})
	}
}

func DelPath(ctx *gin.Context) {
	domainId := ctx.Param("domain_id")
	pathId := ctx.Param("path_id")
	if domainId == "" || pathId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	deleted := delPathDataByDomainIdAndPathId(domainId, pathId)
	if deleted {
		base.Result{Context: ctx}.SucResult(nil)
		return
	}
	base.Result{Context: ctx}.ErrResult(base.SystemError)
}
