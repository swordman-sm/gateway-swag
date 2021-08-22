package impl

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/handler"
	"github.com/form3tech-oss/jwt-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type AuthServiceImpl struct {
}

//根据userId获取存储在etcd指定的user key
func getUserKey(userId string) string {
	return fmt.Sprintf(base.AdminUserDataPathFormat, userId)
}

//在etcd中根据用户id获取用户
func (AuthServiceImpl) GetAdminByUserId(userId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	resp, err := base.Cli.Get(ctx, getUserKey(userId))
	cancel()
	return resp, err
}

//初始化admin用户信息
func (AuthServiceImpl) InitAdminData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, base.AuthInitDataPath)
	cancel()
	return rsp, err
}

//用户名密码MD5加密
func (AuthServiceImpl) Md5UsernameAndPwd(username, password, salt string) (string, string) {
	//用户名加密
	digest := md5.New()
	digest.Write([]byte(username))
	md5UserId := digest.Sum(nil)
	userIdHexStr := hex.EncodeToString(md5UserId)
	//密码加密
	passwordHexStr := ""
	if password != "" && salt != "" {
		digest.Reset()
		digest.Write([]byte(salt))
		digest.Write([]byte(password))
		digest.Write([]byte(salt))
		md5Pwd := digest.Sum(nil)
		passwordHexStr = hex.EncodeToString(md5Pwd)
	}
	return userIdHexStr, passwordHexStr
}

func (AuthServiceImpl) AddNewAdmin(userId string, adminJson []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	//etcd 使用 Txn 提供简单的事务处理，使用这个特性，可以一次性插入多条语句
	cancel()
	txn := base.Cli.Txn(ctx)
	commit, err := txn.Then(clientv3.OpPut(getUserKey(userId), string(adminJson)),
		clientv3.OpPut(base.AuthDataPath, time.Now().Format("2006-01-02 15:04"))).Commit()
	if err != nil {
		return false
	}
	return commit.Succeeded
}

//获取token jwt加密token
func (AuthServiceImpl) GetToken(userId, salt string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(handler.TokenExpire).Unix(),
		Issuer:    handler.Issuser,
		Subject:   userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(salt))
	if err != nil {
		return ""
	}
	return tokenStr
}

//check token是否有效
func (AuthServiceImpl) CheckToken(userId, salt, jwtStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})
	return token, err
}
