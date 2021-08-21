package modules

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func gatewayActivePath(serName string) string {
	return fmt.Sprintf(base.MetricsGatewayActiveFormat, serName)
}

func gatewayActivesDataPath(serName string) string {
	return fmt.Sprintf(base.MetricsGatewayActivesDataFormat, serName)
}

func gatewaysData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, base.MetricsGatewayActivePrefix, clientv3.WithPrefix())
	cancel()
	return rsp, err
}

func gatewayMachineData(serName string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, gatewayActivePath(serName))
	cancel()
	return rsp, err
}

func gatewayData(serName string, limit int64) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	op := clientv3.WithLastKey()
	op = append(op, clientv3.WithLimit(limit))
	rsp, err := base.Cli.Get(ctx, gatewayActivesDataPath(serName), op...)
	cancel()
	return rsp, err
}
