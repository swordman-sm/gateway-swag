package domain

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var sysLog *logrus.Logger
var proxyLog *logrus.Logger

func init() {
	sysLog = logrus.New()
	proxyLog = logrus.New()
	proxyLog.Formatter = &logrus.JSONFormatter{}
}

func setNull(logger *logrus.Logger) {
	//进程从null device中读取（os.DevNull），stdin也可以时一个文件，否则的话则在运行过程中再开一个goroutine去
	fileWriter, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logger.Out = fileWriter
}

func Sys() *logrus.Logger {
	return sysLog
}

func Proxy() *logrus.Logger {
	return proxyLog
}
