package modules

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//域名数据的key
func getDomainDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainDataFormat, domainId)
}

//域名备份数据的key
func getDomainBakDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainBakDataFormat, domainId)
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

//获取所有域名数据
func getAllDomainsData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	response, err := base.Cli.Get(ctx, base.DomainsDataPrefix, clientv3.WithPrefix())
	cancel()
	return response, err
}

//获取指定id域名数据
func getDomainDataByDomainId(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, getDomainDataKey(domainId))
	cancel()
	return rsp, err
}

//存储域名数据
func addDomainData(domainId, domainJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	_, err := base.Cli.Put(ctx, getDomainDataKey(domainId), domainJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除指定domainId数据并备份数据
func delDomainDataByDomainId(domainId string) bool {
	dataKey := getDomainDataKey(domainId)
	dataBakKey := getDomainBakDataKey(domainId)
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	defer cancel()
	dataRsp, err := base.Cli.Get(ctx, dataKey)
	if err != nil {
		return false
	}
	data := dataRsp.Kvs[0].Value

	txn := base.Cli.Txn(ctx)
	lease, err := base.Cli.Grant(ctx, base.BakDataTTL)
	if err != nil {
		return false
	}
	rsp, err := txn.Then(clientv3.OpDelete(dataKey),
		clientv3.OpPut(dataBakKey, string(data), clientv3.WithLease(lease.ID))).Commit()
	if err != nil {
		return false
	}
	return rsp.Succeeded
}

//指定domainId所有路径数据
func getPathDataByDomainId(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, getDomainPathsDataKey(domainId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//指定domainId的pathId路径数据
func getPathDataByDomainIdAndPathId(domainId string, pathId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, getDomainPathDataKey(domainId, pathId), clientv3.WithPrefix())
	cancel()
	return rsp, err
}

//设置路径数据
func addPathData(domainId, pathId string, pathJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	_, err := base.Cli.Put(ctx, getDomainPathDataKey(domainId, pathId), pathJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除指定domainId和pathId的路径数据并备份数据
func delPathDataByDomainIdAndPathId(domainId string, pathId string) bool {
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
