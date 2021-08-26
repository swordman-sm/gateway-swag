package impl

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"gateway-swag/gateway/base"
	"gateway-swag/gateway/domain"
	"gateway-swag/gateway/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
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
				base.Sys().Warnf("证书生成失败 %s", string(cert.SerName))
			}
			certMap[cert.SerName] = &certificate
		}
		base.Sys().Infoln("所有域名证书设置完成")
	}
}

func (CertServiceImpl) GetCertData(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
	if certMap == nil {
		return nil, errors.New("暂未发现有效证书信息")
	}
	lowerServerName := strings.ToLower(info.ServerName)
	if cert, ok := certMap[lowerServerName]; ok {
		return cert, nil
	}
	tags := strings.Split(lowerServerName, ".")
	for i := range tags {
		//通配路径证书查询
		tags[i] = "*"
		tmpKey := strings.Join(tags, ".")
		if cert, ok := certMap[tmpKey]; ok {
			return cert, nil
		}
	}
	return nil, errors.New("find no certificate")
}

func (CertServiceImpl) CertDataChangeListen() {
	//NIO 监听证书变化情况
	ech := make(chan *clientv3.Event, 100)
	go etcd.WatchCertChange(ech)
	for {
		select {
		case event := <-ech:
			if event.Type == clientv3.EventTypePut {
				cert := new(domain.Cert)
				err := json.Unmarshal(event.Kv.Value, cert)
				if err != nil {
					base.Sys().Warnf("证书数据解析失败 %s", string(event.Kv.Value))
					continue
				}
				certificate, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock))
				if err != nil {
					base.Sys().Warnf("证书生成失败 %s", string(event.Kv.Value))
				}
				certMap[cert.SerName] = &certificate
				base.Sys().Infof("域名%s证书更新完成", cert.SerName)
			} else if event.Type == clientv3.EventTypeDelete {
				//删除操作后已将数据备份至bak中
				certBakData, err := etcd.GetCertDataByPath(etcd.GetCertBakDataKey(string(event.Kv.Key)))
				//未找到备份数据->失败
				if err != nil {
					base.Sys().Warnf("【域名证书路径%s】删除-获取备份数据失败", string(event.Kv.Key))
					continue
				}
				delete(certMap, certBakData.SerName)
			}
		}
	}
}
