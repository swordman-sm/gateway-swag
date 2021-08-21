package modules

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var requestListenTTL int64 = 60

func requestListenDataK(listenId string) string {
	return fmt.Sprintf(base.RequestListenDataFormat, listenId)
}
func putRequestListen(listenId string, copyJson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	lease, err := base.Cli.Grant(ctx, requestListenTTL)
	if err != nil {
		return err
	}
	_, err = base.Cli.Put(ctx, requestListenDataK(listenId), copyJson, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func requestsCopy() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	op := clientv3.WithLastKey()
	op = append(op, clientv3.WithLimit(500))
	rsp, err := base.Cli.Get(ctx, base.RequestsCopyDataPrefix, op...)
	cancel()
	return rsp, err
}
