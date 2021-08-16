package main

import (
	"gateway-swag/management/modules"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

)

var (
	//管理后台service启动地址
	addr = kingpin.Flag("addr", "gateway listen addr").Required().String()
	//连接的etcd地址
	etcd = kingpin.Flag("etcd", "etcd server addr").Required().String()
	//etcd用户名
	username = kingpin.Flag("u", "etcd username").Default("").String()
	//etcd密码
	password = kingpin.Flag("p", "etcd password").Default("").String()
)

//gateway管理后台UI后端逻辑
func main() {
	//设定help简化为-h
	kingpin.HelpFlag.Short('h')
	//解析命令行传递参数
	kingpin.Parse()

	//我们每次启动gin服务器，如果不加
	//gin.SetMode(gin.ReleaseMode)
	//这一段，就会输出一段提示
	//[WARNING] Running in "debug" mode. Switch to "release" mode in production.
	// - using env:   export GIN_MODE=release
	// - using code:  gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	//静态文件
	engine.Static("/admin/", "webapp")
	engine.Static("/static/", "webapp/static")

	//统一群组 路径定义
	engine.Group("/v1", modules.AuthHandler)
	engine.GET("/index", modules.Index, modules.AuthHandler)
	engine.POST("/init", modules.AuthInit)
	engine.POST("/login", modules.Login)
	engine.GET("/logout", modules.Logout)

	logrus.Infof("Gateway 后端管理服务启动地址: %s, etcd服务地址: %s", *addr, *etcd)
	err := engine.Run(*addr)
	if err != nil {
		logrus.Errorf("启动服务失败, 端口监听 %v", err)
	}

}
