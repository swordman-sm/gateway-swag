package service

import clientv3 "go.etcd.io/etcd/client/v3"

type GatewayService interface {
	GetAllGatewaysData() (*clientv3.GetResponse, error)
	GetGatewayDataByServer(serName string) (*clientv3.GetResponse, error)
	GetGatewayDataByLimit(serName string, limit int64) (*clientv3.GetResponse, error)
}
