package modules

import (
	"encoding/json"
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
	Issuser = "swagw_admin"
	//token 过期时间
	TokenExpire = time.Hour * 24
)

//管理员结构体
type AdminUser struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

//不使用的变量可以使用_代替
func Index(_ *gin.Context) {

}

func AuthHandler(ctx *gin.Context) {
	jwtStr, _ := ctx.Cookie("jwt")
	userId, _ := ctx.Cookie("userId")
	//cookie中找不到对应的需求cookie值
	if jwtStr == "" || userId == "" {
		resultCtx{ctx}.ErrResult(SystemErrorNotLogin)
		return
	}
	//根据userId在etcd中找不到对应的用户信息
	userRsp, err := getAdminUserByUserId(userId)
	if err != nil || userRsp.Count == 0 {
		resultCtx{ctx}.ErrResult(SystemErrorNotLogin)
		return
	}
	admin := new(AdminUser)
	//解析从etcd获取的json存储的admin用户信息
	err = json.Unmarshal(userRsp.Kvs[0].Value, admin)
	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	//通过盐值解析jwt声明token信息
	token, err := checkToken(userId, admin.Salt, jwtStr)

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
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
}

func AuthInit(ctx *gin.Context) {
	resp, err := authDataInit()
	if err != nil || resp.Count > 0 {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	//创建一个salt值 这里选择uuid作为产生salt的来源
	salt := uuid.Must(uuid.NewV4()).String()

	//从web获取提交的用户名和密码
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		logrus.Error("username or password is nil")
		return
	}

	userId := initUserNameByMd5(username)
	//盐值增益加密
	md5pwd := initPasswordByMd5(password, salt)

	//创建一个AdminUser用户信息
	adminUser := AdminUser{
		UserId:   userId,
		UserName: username,
		Password: md5pwd,
		Salt:     salt,
	}
	adminJson, err := json.Marshal(adminUser)
	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	//保存创建的AdminUser信息
	isSuccuss := putAdminUser(userId, adminJson)
	if !isSuccuss {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	resultCtx{ctx}.SucResult(SystemSuccess)
}

func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		logrus.Error("username or password is nil")
		return
	}
	//根据用户名获取userId
	userId := initUserNameByMd5(username)
	//通过userId获取存储在etcd中的userJson
	resp, err := getAdminUserByUserId(userId)
	if err != nil || resp.Count == 0 {
		resultCtx{ctx}.ErrResult(LoginParamsError)
		return
	}
	adminUser := new(AdminUser)
	err = json.Unmarshal(resp.Kvs[0].Value, adminUser)

	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	//对输入密码按照md5加密方式加密后与存储在etcd中用户密码信息对比
	if adminUser.Password != initPasswordByMd5(password, adminUser.Salt) {
		resultCtx{ctx}.ErrResult(LoginParamsError)
		return
	}

	//验证通过,获取token
	token := getToken(userId, adminUser.Salt)
	if token == "" {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	ctx.SetCookie("jwt", token, int(TokenExpire.Seconds()), "/", "", false, false)
	ctx.SetCookie("userId", userId, int(TokenExpire.Seconds()), "/", "", false, false)
	resultCtx{ctx}.SucResult(nil)
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "/", "", false, true)
	ctx.SetCookie("userId", "", -1, "/", "", false, true)
	resultCtx{ctx}.SucResult(nil)
}
