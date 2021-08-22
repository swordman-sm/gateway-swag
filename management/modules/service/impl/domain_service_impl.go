package impl

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DomainServiceImpl struct {
}

//域名数据的key
func getDomainDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainDataFormat, domainId)
}

//域名备份数据的key
func getDomainBakDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainBakDataFormat, domainId)
}

//获取所有域名数据
func (DomainServiceImpl) GetAllDomainsData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	response, err := base.Cli.Get(ctx, base.DomainsDataPrefix, clientv3.WithPrefix())
	cancel()
	return response, err
}

//获取指定id域名数据
func (DomainServiceImpl) GetDomainDataByDomainId(domainId string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, getDomainDataKey(domainId))
	cancel()
	return rsp, err
}

//存储域名数据
func (DomainServiceImpl) AddDomainData(domainId, domainJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	_, err := base.Cli.Put(ctx, getDomainDataKey(domainId), domainJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

//删除指定domainId数据并备份数据
func (DomainServiceImpl) DelDomainDataByDomainId(domainId string) bool {
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
