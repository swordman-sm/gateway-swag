package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway-swag/gateway/base"
	"gateway-swag/gateway/domain"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

func GetCertDataKey(certId string) string {
	return fmt.Sprintf(base.SwagCertFormat, certId)
}

func GetCertBakDataKey(certPath string) string {
	return strings.Replace(certPath, base.SwagCertsPrefix, base.SwagCertsBakPrefix, 1)
}

//监听证书数据变化
func WatchCertChange(e chan *clientv3.Event) {
	for {
		watchChan := base.Cli.Watch(context.Background(), base.SwagCertsPrefix, clientv3.WithPrefix())
		for resp := range watchChan {
			for _, event := range resp.Events {
				e <- event
			}
		}
	}
}

//获取所有现有证书数据
func GetAllCertsData() ([]*domain.Cert, error) {
	var certs []*domain.Cert
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	resp, err := base.Cli.Get(ctx, base.SwagCertsPrefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return certs, err
	}
	for _, kv := range resp.Kvs {
		cert := new(domain.Cert)
		err := json.Unmarshal(kv.Value, cert)
		if err == nil {
			certs = append(certs, cert)
		} else {
			base.Sys().Warnf("证书数据解析失败 %s, error: %s", string(kv.Value), err)
		}
	}
	return certs, err
}

//获取指定路径的证书
func GetCertDataByPath(certPath string) (*domain.Cert, error) {
	cert := new(domain.Cert)
	ctx, cancel := context.WithTimeout(context.Background(), base.WriteTimeout)
	resp, err := base.Cli.Get(ctx, certPath, clientv3.WithPrefix())
	cancel()
	if err != nil || resp.Count == 0 {
		return cert, err
	}
	err = json.Unmarshal(resp.Kvs[0].Value, cert)
	if err != nil {
		base.Sys().Warnf("[%s] json parse error, value : %q err: %s", certPath, resp.Kvs[0].Value, err)
		return cert, err
	}
	return cert, err
}
