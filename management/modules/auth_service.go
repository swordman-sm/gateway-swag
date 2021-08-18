package modules

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//根据userId获取存储在etcd指定的user key
func getUserKey(userId string) string {
	return fmt.Sprintf(adminUserDataPathFormat, userId)
}

//在etcd中根据用户id获取用户
func getAdminUserByUserId(userId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	resp, err := cli.Get(ctx, getUserKey(userId))
	cancel()
	return resp, err
}

//初始化admin用户信息
func authDataInit() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	rsp, err := cli.Get(ctx, authInitDataPath)
	cancel()
	return rsp, err
}

//将用户名转16进制MD5加密字符串
func initUserNameByMd5(username string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(username))
	md5UserName := md5Ctx.Sum(nil)
	return hex.EncodeToString(md5UserName)
}

//将密码转16进制MD5及盐值加密后的字符串
func initPasswordByMd5(password, salt string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(salt))
	md5Ctx.Write([]byte(password))
	md5Ctx.Write([]byte(salt))
	md5Pwd := md5Ctx.Sum(nil)
	return hex.EncodeToString(md5Pwd)
}

func putAdminUser(userId string, adminJson []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	//etcd 使用 Txn 提供简单的事务处理，使用这个特性，可以一次性插入多条语句
	txn := cli.Txn(ctx)
	commit, err := txn.Then(clientv3.OpPut(getUserKey(userId), string(adminJson)),
		clientv3.OpPut(authDataPath, time.Now().Format("2006-01-02 15:04"))).Commit()
	if err != nil {
		return false
	}
	cancel()
	return commit.Succeeded
}

//check token是否有效
func checkToken(userId, salt, jwtStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})
	return token, err
}

//获取token jwt加密token
func getToken(userId, salt string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(TokenExpire).Unix(),
		Issuer:    Issuser,
		Subject:   userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(salt))
	if err != nil {
		return ""
	}
	return tokenStr
}
