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
	SwagCertsPrefix    = SwagPrefix + "server-tls/"
	SwagCertsBakPrefix = SwagPrefix + "server-tls-bak/"
	SwagCertFormat     = SwagCertsPrefix + "%s"
	SwagCertBakFormat  = SwagCertsBakPrefix + "%s"
	//Domain relate
	DomainsDataPrefix    = SwagPrefix + "domain-data/"
	DomainsBakDataPrefix = SwagPrefix + "domain-data-bak/"
	DomainDataFormat     = DomainsDataPrefix + "%s/"
	DomainBakDataFormat  = DomainsBakDataPrefix + "%s/"
	//Path relate
	DomainPathsDataPrefix    = SwagPrefix + "path-data/"
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
	//iss ??? OpenId Connect???????????????OIDC???????????????????????????????????????????????? ???Issuer Identifier??????
	//??????????????????????????????????????????????????? Token ?????????????????????????????????????????? http(s) url?????? https://www.baidu.com???
	//??? Token ?????????????????????????????????????????????????????????????????????????????????????????????????????????????????? HTTP 401
	Issuser = "swag_admin"
	//token ????????????
	TokenExpire = time.Hour * 24
)
