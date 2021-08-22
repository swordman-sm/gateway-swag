package base

import (
	"time"
)

//Error Level
const (
	//System level
	SystemSuccess       = 200
	SystemError         = -1001
	SystemErrorNotInit  = -2000
	SystemErrorNotLogin = -2001
	//Login level
	LoginParamsError = -2002
	//Parse level
	JsonParseError        = -2003
	DataCannotDeleteError = -3001
	DataParseError        = -4001
)

//Timeout duration
const (
	DialTimeout  = 3 * time.Second
	ReadTimeout  = 3 * time.Second
	WriteTimeout = 3 * time.Second
	BakDataTTL   = 1800
)

//ETCD Path relate
const (
	SwagPrefix = "/swag-gateway/"

	//Auth path
	AuthDataPath            = SwagPrefix + "auth-data"
	AuthInitDataPath        = AuthDataPath + "/init"
	AdminUserDataPathFormat = AuthDataPath + "/user/%s"
	//Cert relate
	HgwCertsPrefix    = SwagPrefix + "server-tls/"
	HgwCertsBakPrefix = SwagPrefix + "server-tls-bak/"
	HgwCertFormat     = HgwCertsPrefix + "%s"
	HgwCertBakFormat  = HgwCertsBakPrefix + "%s"
	//Domain relate
	DomainsDataPrefix    = SwagPrefix + "domain-data/"
	DomainsBakDataPrefix = SwagPrefix + "domain-data-bak/"
	DomainDataFormat     = DomainsDataPrefix + "%s/"
	DomainBakDataFormat  = DomainsBakDataPrefix + "%s/"
	//Path relate
	DomainPathsBakDataPrefix = SwagPrefix + "path-data-bak/%s/"
	DomainPathsDataFormat    = SwagPrefix + "path-data/%s/"
	DomainPathDataFormat     = DomainPathsDataFormat + "%s"
	DomainPathBakDataFormat  = DomainPathsBakDataPrefix + "%s"
	//Gateway metrics relate
	MetricsGatewayActivePrefix      = SwagPrefix + "gateway-active/"
	MetricsGatewayActiveFormat      = MetricsGatewayActivePrefix + "%s"
	MetricsGatewayActiveDataPrefix  = SwagPrefix + "gateway-active-data/"
	MetricsGatewayActivesDataFormat = MetricsGatewayActiveDataPrefix + "%s/"
	MetricsGatewayActiveDataFormat  = MetricsGatewayActivesDataFormat + "%s"

	//RequestsListen relate
	RequestsListenDataPrefix = SwagPrefix + "requests-listen/"
	RequestListenDataFormat  = RequestsListenDataPrefix + "%s"
	//RequestsCopy relate
	RequestsCopyDataPrefix = SwagPrefix + "requests-copy/"
)

//JWT relate
const (
	//iss 是 OpenId Connect（后文简称OIDC）协议中定义的一个字段，其全称为 “Issuer Identifier”，
	//中文意思就是：颁发者身份标识，表示 Token 颁发者的唯一标识，一般是一个 http(s) url，如 https://www.baidu.com。
	//在 Token 的验证过程中，会将它作为验证的一个阶段，如无法匹配将会造成验证失败，最后返回 HTTP 401
	Issuser = "swag_admin"
	//token 过期时间
	TokenExpire = time.Hour * 24
)
