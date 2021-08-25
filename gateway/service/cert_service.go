package service

import "crypto/tls"

type CertService interface {
	//初始化所有证书设置
	InitCertsData()
	GetCertData(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error)
	CertDataChangeListen()
}
