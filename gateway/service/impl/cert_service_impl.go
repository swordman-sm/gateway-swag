package impl

import (
	"crypto/tls"
	"errors"
	"gateway-swag/gateway/domain"
	"gateway-swag/gateway/etcd"
)

type CertServiceImpl struct{}

var certMap map[string]*tls.Certificate

//仅初始化一次
func init() {
	certMap = make(map[string]*tls.Certificate)
}

func (CertServiceImpl) InitCertsData() {
	certs, err := etcd.GetAllCertsData()
	if err == nil {
		for _, cert := range certs {
			//解析证书数据
			certificate, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock))
			if err != nil {
				domain.Sys().Warnf("证书生成失败 %s", string(cert.SerName))
			}
			certMap[cert.SerName] = &certificate
		}
		domain.Sys().Infoln("所有域名证书设置完成")
	}
}

func (CertServiceImpl) GetCertData(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
	if certMap == nil {
		return nil, errors.New("暂未发现有效证书信息")
	}
	info.ServerName
}

func (CertServiceImpl) CertDataChangeListen() {
	panic("implement me")
}
