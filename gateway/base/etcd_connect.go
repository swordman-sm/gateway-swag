package base

import clientv3 "go.etcd.io/etcd/client/v3"

var Cli *clientv3.Client

func EtcdConnect(endPoints []string, username, password string) error {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: DialTimeout,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		return err
	}
	Cli = c
	return nil
}
