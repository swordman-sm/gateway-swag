package impl

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type CertServiceImpl struct {
}

func getCertDataKey(certId string) string {
	return fmt.Sprintf(base.HgwCertFormat, certId)
}

func getCertBakDataKey(certId string) string {
	return fmt.Sprintf(base.HgwCertBakFormat, certId)
}

func (CertServiceImpl) GetAllCertData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	resp, err := base.Cli.Get(ctx, base.HgwCertsPrefix, clientv3.WithPrefix())
	cancel()
	return resp, err
}

func (CertServiceImpl) AddCertData(certId, certJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	_, err := base.Cli.Put(ctx, getCertDataKey(certId), certJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (CertServiceImpl) DelCertData(certId string) bool {
	certDataKey := getCertDataKey(certId)
	certBakDataKey := getCertBakDataKey(certId)
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	defer cancel()
	resp, err := base.Cli.Get(ctx, certDataKey)
	if err != nil {
		return false
	}
	certData := resp.Kvs[0].Value
	//开启事务
	txn := base.Cli.Txn(ctx)
	lease, err := base.Cli.Grant(ctx, base.BakDataTTL)
	if err != nil {
		return false
	}
	res, err := txn.Then(clientv3.OpDelete(certDataKey),
		clientv3.OpPut(certBakDataKey, string(certData), clientv3.WithLease(lease.ID)),
	).Commit()
	if err != nil {
		return false
	}
	return res.Succeeded
}
