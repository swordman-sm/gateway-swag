package service

import clientv3 "go.etcd.io/etcd/client/v3"

type RequestListenService interface {
	AddRequestListen(listenId string, copyJson string) error
	GetRequestsCopy() (*clientv3.GetResponse, error)
}
