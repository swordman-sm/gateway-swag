package service

import clientv3 "go.etcd.io/etcd/client/v3"

type PathService interface {
	//Path Data
	//获取指定domainId所有路径数据
	GetPathsDataByDomainId(domainId string) (*clientv3.GetResponse, error)
	//获取指定domainId的pathId路径数据
	GetPathDataByDomainIdAndPathId(domainId string, pathId string) (*clientv3.GetResponse, error)
	//添加路径数据
	AddPathData(domainId, pathId string, pathJson string) error
	//删除指定domainId和pathId的路径数据并备份数据
	DelPathDataByDomainIdAndPathId(domainId string, pathId string) bool
}
