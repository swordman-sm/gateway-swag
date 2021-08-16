package modules

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	SystemSuccess           = 200
	SystemError             = -1001
	SystemErrorNotInit      = -2000
	SystemErrorNotLogin     = -2001
	LoginParamsError        = -2002
	DataCannotDeleteError   = -3001
	DataParseError          = -4001
	hgwPrefix               = "/swag-gateway/"
	authDataPath            = hgwPrefix + "auth-data/"
	authInitDataPath        = authDataPath + "/init"
	adminUserDataPathFormat = authDataPath + "/user/%s"
	dialTimeout             = 3 * time.Second
	readTimeout             = 3 * time.Second
	writeTimeout            = 3 * time.Second
	bakDataTTL              = 1800
)

//成功信息
type SuccessOutPut struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

//错误信息
type ErrorOutPut struct {
	Status    int `json:"status"`
	ErrorCode int `json:"error_code"`
}

//包装ctx 增加api func操作
type resultCtx struct {
	*gin.Context
}

func (r resultCtx) SucResult(data interface{}) {
	r.JSON(http.StatusOK, SuccessOutPut{Status: 1, Data: data})
	r.Next()
}

func (r resultCtx) ErrResult(errorCode int) {
	r.AbortWithStatusJSON(http.StatusOK, ErrorOutPut{Status: 0, ErrorCode: errorCode})
}
