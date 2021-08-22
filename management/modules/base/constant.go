package base

import (
	"time"
)

const (
	SystemSuccess           = 200
	SystemError             = -1001
	SystemErrorNotInit      = -2000
	SystemErrorNotLogin     = -2001
	LoginParamsError        = -2002
	DataCannotDeleteError   = -3001
	DataParseError          = -4001
	HgwPrefix               = "/swag-gateway/"
	AuthDataPath            = HgwPrefix + "auth-data"
	AuthInitDataPath        = AuthDataPath + "/init"
	AdminUserDataPathFormat = AuthDataPath + "/user/%s"
	DialTimeout             = 3 * time.Second
	ReadTimeout             = 3 * time.Second
	WriteTimeout            = 3 * time.Second
	BakDataTTL              = 1800
	HgwCertsPrefix          = HgwPrefix + "server-tls/"
	HgwCertFormat           = HgwCertsPrefix + "%s"
	HgwCertsBakPrefix       = HgwPrefix + "server-tls-bak/"
	HgwCertBakFormat        = HgwCertsBakPrefix + "%s"
	DomainsDataPrefix       = HgwPrefix + "domain-data/"
	DomainDataFormat        = DomainsDataPrefix + "%s/"
	DomainPathsDataFormat   = HgwPrefix + "path-data/%s/"
	DomainPathDataFormat    = DomainPathsDataFormat + "%s"

	DomainsBakDataPrefix = HgwPrefix + "domain-data-bak/"
	DomainBakDataFormat  = DomainsBakDataPrefix + "%s/"

	DomainPathsBakDataPrefix = HgwPrefix + "path-data-bak/%s/"
	DomainPathBakDataFormat  = DomainPathsBakDataPrefix + "%s"

	MetricsGatewayActivePrefix      = HgwPrefix + "gateway-active/"
	MetricsGatewayActiveFormat      = MetricsGatewayActivePrefix + "%s"
	MetricsGatewayActiveDataPrefix  = HgwPrefix + "gateway-active-data/"
	MetricsGatewayActivesDataFormat = MetricsGatewayActiveDataPrefix + "%s/"
	MetricsGatewayActiveDataFormat  = MetricsGatewayActivesDataFormat + "%s"

	RequestsListenDataPrefix = HgwPrefix + "requests-listen/"
	RequestListenDataFormat  = RequestsListenDataPrefix + "%s"

	RequestsCopyDataPrefix = HgwPrefix + "requests-copy/"
)
