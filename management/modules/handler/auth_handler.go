package handler

import (
	"encoding/json"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"gateway-swag/management/modules/service/impl"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	//iss 是 OpenId Connect（后文简称OIDC）协议中定义的一个字段，其全称为 “Issuer Identifier”，
	//中文意思就是：颁发者身份标识，表示 Token 颁发者的唯一标识，一般是一个 http(s) url，如 https://www.baidu.com。
	//在 Token 的验证过程中，会将它作为验证的一个阶段，如无法匹配将会造成验证失败，最后返回 HTTP 401
	Issuser = "swag_admin"
	//token 过期时间
	TokenExpire = time.Hour * 24
)

var authService = new(impl.AuthServiceImpl)

//不使用的变量可以使用_代替
func IndexHandler(_ *gin.Context) {

}

func AuthHandler(ctx *gin.Context) {
	jwtStr, _ := ctx.Cookie("jwt")
	userId, _ := ctx.Cookie("userId")
	//cookie中找不到对应的需求cookie值
	if jwtStr == "" || userId == "" {
		base.Result{Context: ctx}.ErrResult(base.SystemErrorNotLogin)
		return
	}
	//根据userId在etcd中找不到对应的用户信息

	userRsp, err := authService.GetAdminByUserId(userId)
	if err != nil || userRsp.Count == 0 {
		base.Result{Context: ctx}.ErrResult(base.SystemErrorNotLogin)
		return
	}
	admin := new(domain.AdminUser)
	//解析从etcd获取的json存储的admin用户信息
	err = json.Unmarshal(userRsp.Kvs[0].Value, admin)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	//通过盐值解析jwt声明token信息
	token, err := authService.CheckToken(userId, admin.Salt, jwtStr)

	if token != nil && err == nil {
		//token验证通过
		if token.Valid {
			//存放验证信息及其时间信息 类型转换成标准
			claims := token.Claims.(*jwt.StandardClaims)
			if userId == claims.Subject {
				//执行其他中间件
				ctx.Next()
				return
			}
		}
	}
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
}

func InitAdminHandler(ctx *gin.Context) {
	resp, err := authService.InitAdminData()
	if err != nil || resp.Count > 0 {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	//创建一个salt值 这里选择uuid作为产生salt的来源
	salt := uuid.Must(uuid.NewV4()).String()

	//从web获取提交的用户名和密码
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		logrus.Error("username or password is nil")
		return
	}

	userId, md5pwd := authService.Md5UsernameAndPwd(username, password, salt)

	//创建一个AdminUser用户信息
	adminUser := domain.AdminUser{
		UserId:   userId,
		UserName: username,
		Password: md5pwd,
		Salt:     salt,
	}
	adminJson, err := json.Marshal(adminUser)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	//保存创建的AdminUser信息
	isSuccuss := authService.AddNewAdmin(userId, adminJson)
	if !isSuccuss {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	base.Result{Context: ctx}.SucResult(base.SystemSuccess)
}

func LoginHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		logrus.Error("username or password is nil")
		return
	}
	//根据用户名获取userId
	userId, _ := authService.Md5UsernameAndPwd(username, "", "")
	//通过userId获取存储在etcd中的userJson
	resp, err := authService.GetAdminByUserId(userId)
	if err != nil || resp.Count == 0 {
		base.Result{Context: ctx}.ErrResult(base.LoginParamsError)
		return
	}

	adminUser := new(domain.AdminUser)
	err = json.Unmarshal(resp.Kvs[0].Value, adminUser)

	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	//对输入密码按照md5加密方式加密后与存储在etcd中用户密码信息对比
	_, md5Pwd := authService.Md5UsernameAndPwd(username, password, adminUser.Salt)
	if adminUser.Password != md5Pwd {
		base.Result{Context: ctx}.ErrResult(base.LoginParamsError)
		return
	}

	//验证通过,获取token
	token := authService.GetToken(userId, adminUser.Salt)
	if token == "" {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	ctx.SetCookie("jwt", token, int(TokenExpire.Seconds()), "/", "", false, false)
	ctx.SetCookie("userId", userId, int(TokenExpire.Seconds()), "/", "", false, false)
	base.Result{Context: ctx}.SucResult("success")
}

func LogoutHandler(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "/", "", false, true)
	ctx.SetCookie("userId", "", -1, "/", "", false, true)
	base.Result{Context: ctx}.SucResult("success")
}
