package service

import clientv3 "go.etcd.io/etcd/client/v3"

type DomainService interface {
	//Domain Data
	//获取所有域名数据
	GetAllDomainsData() (*clientv3.GetResponse, error)
	//获取指定id域名数据
	GetDomainDataByDomainId(domainId string) (*clientv3.GetResponse, error)
	//存储域名数据
	AddDomainData(domainId, domainJson string) error
	//删除指定domainId数据并备份数据
	DelDomainDataByDomainId(domainId string) bool
}
