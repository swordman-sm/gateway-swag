package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//成功信息
type SuccessData struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

//错误信息
type ErrorData struct {
	Status    int `json:"status"`
	ErrorCode int `json:"error_code"`
}

//包装ctx 增加api func操作
type Result struct {
	Context *gin.Context
}

func (r Result) SucResult(data interface{}) {
	r.Context.JSON(http.StatusOK, SuccessData{Status: 1, Data: data})
	r.Context.Next()
}

func (r Result) ErrResult(errorCode int) {
	r.Context.AbortWithStatusJSON(http.StatusOK, ErrorData{Status: 0, ErrorCode: errorCode})
}
