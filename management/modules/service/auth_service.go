package service

import (
	"github.com/form3tech-oss/jwt-go"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type AuthService interface {
	//从etcd中获取Admin的信息
	GetAdminByUserId(userId string) (*clientv3.GetResponse, error)
	//初始化Admin数据
	InitAdminData() (*clientv3.GetResponse, error)
	//用户名/密码MD5加密
	Md5UsernameAndPwd(username, password, salt string) (string, string)
	//添加Admin用户
	AddNewAdmin(userId string, adminJson []byte) bool
	//获取token
	GetToken(userId, salt string) string
	//校验token
	CheckToken(userId, salt, jwtStr string) (*jwt.Token, error)
}
