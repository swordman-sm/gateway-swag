package service

import clientv3 "go.etcd.io/etcd/client/v3"

type CertService interface {
	//从etcd获取所有的证书数据
	GetAllCertData() (*clientv3.GetResponse, error)
	//新增或更新证书信息
	AddCertData(certId, certJson string) error
	//删除证书信息
	DelCertData(certId string) bool
}
