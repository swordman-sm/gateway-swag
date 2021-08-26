package main

import (
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/handler"
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

	err := base.EtcdConnect([]string{*etcd}, *username, *password)
	if err != nil {
		panic(err)
	}

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

	engine.GET("/index", handler.IndexHandler, handler.AuthHandler)
	engine.POST("/init", handler.InitAdminHandler)
	engine.POST("/login", handler.LoginHandler)
	engine.GET("/logout", handler.LogoutHandler)

	//统一群组 路径定义
	group := engine.Group("/swag", handler.AuthHandler)

	group.POST("/domains/", handler.AddDomainDataHandler)
	group.GET("/domains/", handler.QueryAllDomainsDataHandler)
	group.POST("/domains/:domain_id", handler.AddDomainDataHandler)
	group.GET("/domains/:domain_id", handler.QueryDomainDataByDomainIdHandler)
	group.DELETE("/domains/:domain_id", handler.DelDomainByDomainIdHandler)

	group.POST("/domains/:domain_id/paths/", handler.AddPathDataHandler)
	group.POST("/domains/:domain_id/paths/:path_id", handler.AddPathDataHandler)
	group.GET("/domains/:domain_id/paths/", handler.GetPathsDataByDomainIdHandler)
	group.GET("/domains/:domain_id/paths/:path_id", handler.GetPathDataByDomainIdHandler)
	group.DELETE("/domains/:domain_id/paths/:path_id", handler.DelPathDataByDomainIdAndPathIdHandler)

	engine.POST("/certs/", handler.AddCertHandler)
	engine.GET("/certs/", handler.GetAllCertsDataHandler)
	engine.POST("/certs/:cert_id", handler.AddCertHandler)
	engine.DELETE("/certs/:cert_id", handler.DelCertByCertIdHandler)

	group.GET("/gateways/", handler.GetAllGatewayDataHandler)
	group.GET("/gateways/:server_name", handler.GetGatewayDataByServerHandler)

	group.POST("/requests-listen/:domain_id/", handler.AddRequestListenHandler)
	group.GET("/requests-copy/", handler.RequestsCopyHandler)

	logrus.Infof("Gateway 后端管理服务启动地址: %s, etcd服务地址: %s", *addr, *etcd)
	err = engine.Run(*addr)
	if err != nil {
		logrus.Errorf("启动服务失败, 端口监听 %v", err)
	}

}
