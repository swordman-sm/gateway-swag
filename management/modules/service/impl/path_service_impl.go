package impl

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PathServiceImpl struct {
}

//域名具体路径数据的Key
func getDomainPathDataKey(domainId string, pathId string) string {
	return fmt.Sprintf(base.DomainPathDataFormat, domainId, pathId)
}

//域名备份具体路径数据的Key
func getDomainPathBakDataKey(domainId string, pathId string) string {
	return fmt.Sprintf(base.DomainPathBakDataFormat, domainId, pathId)
}

//指定域名的路径数据key
func getDomainPathsDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainPathsDataFormat, domainId)
}

//指定domainId所有路径数据
func (PathServiceImpl) GetPathsDataByDomainId(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, getDomainPathsDataKey(domainId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//指定domainId的pathId路径数据
func (PathServiceImpl) GetPathDataByDomainIdAndPathId(domainId string, pathId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, getDomainPathDataKey(domainId, pathId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//设置路径数据
func (PathServiceImpl) AddPathData(domainId, pathId string, pathJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	_, err := base.Cli.Put(ctx, getDomainPathDataKey(domainId, pathId), pathJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除指定domainId和pathId的路径数据并备份数据
func (PathServiceImpl) DelPathDataByDomainIdAndPathId(domainId string, pathId string) bool {
	dataK := getDomainPathDataKey(domainId, pathId)
	dataBakK := getDomainPathBakDataKey(domainId, pathId)
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	defer cancel()
	dataRsp, err := base.Cli.Get(ctx, dataK)
	if err != nil {
		return false
	}
	data := dataRsp.Kvs[0].Value

	txn := base.Cli.Txn(ctx)
	lease, err := base.Cli.Grant(ctx, base.BakDataTTL)
	if err != nil {
		return false
	}
	rsp, err := txn.Then(clientv3.OpDelete(dataK),
		clientv3.OpPut(dataBakK, string(data), clientv3.WithLease(lease.ID))).Commit()
	if err != nil {
		return false
	}
	return rsp.Succeeded
}
