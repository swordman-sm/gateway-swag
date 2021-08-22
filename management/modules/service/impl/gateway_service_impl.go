package impl

import (
	"context"
	"fmt"
	"gateway-swag/management/modules/base"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type GatewayServiceImpl struct {
}

func gatewayActivePath(serName string) string {
	return fmt.Sprintf(base.MetricsGatewayActiveFormat, serName)
}

func gatewayActivesDataPath(serName string) string {
	return fmt.Sprintf(base.MetricsGatewayActivesDataFormat, serName)
}

func (GatewayServiceImpl) GetAllGatewaysData() (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, base.MetricsGatewayActivePrefix, clientv3.WithPrefix())
	cancel()
	return rsp, err
}

func (GatewayServiceImpl) GetGatewayDataByServer(serName string) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	rsp, err := base.Cli.Get(ctx, gatewayActivePath(serName))
	cancel()
	return rsp, err
}

func (GatewayServiceImpl) GetGatewayDataByLimit(serName string, limit int64) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	op := clientv3.WithLastKey()
	//clientv3.WithLimit()来实现分页获取的效果
	op = append(op, clientv3.WithLimit(limit))
	rsp, err := base.Cli.Get(ctx, gatewayActivesDataPath(serName), op...)
	cancel()
	return rsp, err
}
