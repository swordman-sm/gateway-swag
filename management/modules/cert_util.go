package modules

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func getCertDataKey(certId string) string {
	return fmt.Sprintf(hgwCertFormat, certId)
}

func getCertBakDataKey(certId string) string {
	return fmt.Sprintf(hgwCertBakFormat, certId)
}

func getAllCertData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	resp, err := cli.Get(ctx, hgwCertsPrefix, clientv3.WithPrefix())
	cancel()
	return resp, err
}

func putCertData(certId, certJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	_, err := cli.Put(ctx, getCertDataKey(certId), certJson)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func delCertData(certId string) bool {
	certDataKey := getCertDataKey(certId)
	certBakDataKey := getCertBakDataKey(certId)
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	resp, err := cli.Get(ctx, certDataKey)
	if err != nil {
		return false
	}
	certData := resp.Kvs[0].Value
	//开启事务
	txn := cli.Txn(ctx)
	lease, err := cli.Grant(ctx, bakDataTTL)
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
