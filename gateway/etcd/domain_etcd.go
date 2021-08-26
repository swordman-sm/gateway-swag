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

//域名数据key
func getDomainDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainDataFormat, domainId)
}

//获取域名路径数据备份路径 替换正常数据路径
func GetDomainPathBakDataKey(dataKey string) string {
	return strings.Replace(dataKey, base.DomainPathsDataPrefix, base.DomainPathsBakDataPrefix, 1)
}

//域名key装域名备份的key
func GetDomainDataKey(domainKey string) string {
	return strings.Replace(domainKey, base.DomainsDataPrefix, base.DomainsBakDataPrefix, 1)
}

//域名下path数据的路径
func getDomainPathsDataKey(domainId string) string {
	return fmt.Sprintf(base.DomainPathsDataFormat, domainId)
}

func GetDomainDataByKey(domainKey string, getPaths bool) (*domain.Domain, error) {
	return getDomainDataByKey(domainKey, getPaths)
}

func GetDomainDataById(domainId string, getPaths bool) (*domain.Domain, error) {
	return getDomainDataByKey(getDomainDataKey(domainId), getPaths)
}

//获取所有域名定义数据
func getAllDomainsData() ([]*domain.Domain, error) {
	return allDomainsData(base.DomainsDataPrefix)
}

//获取域名定义数据
func getDomainDataByKey(dataKey string, getPaths bool) (*domain.Domain, error) {
	domainD := new(domain.Domain)
	ctx, cancel := context.WithTimeout(context.Background(), base.ReadTimeout)
	rsp, err := base.Cli.Get(ctx, dataKey, clientv3.WithPrefix())
	cancel()
	if err != nil || rsp.Count == 0 {
		return domainD, err
	}
	err = json.Unmarshal(rsp.Kvs[0].Value, domainD)
	if err != nil {
		base.Sys().Warnf("域名数据json解析失败 key : %s val: %s err: %s", dataKey, string(rsp.Kvs[0].Value), err)
	} else {
		if getPaths {
			pathsData, err := allDomainPathsData(getDomainPathsDataKey(domainD.Id))
			if err != nil {
				domainD.Paths = pathsData
			}
		}
	}
	return domainD, nil
}

//获取所有域名定义数据
func allDomainsData(dataKey string) ([]*domain.Domain, error) {
	var domainsData []*domain.Domain
	ctx, cancel := context.WithTimeout(context.Background(), base.ReadTimeout)
	rsp, err := base.Cli.Get(ctx, dataKey, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return domainsData, err
	}
	for _, rr := range rsp.Kvs {
		domainD := new(domain.Domain)
		err := json.Unmarshal(rr.Value, domainD)
		if err != nil {
			base.Sys().Warnf("域名数据json解析失败 key : %s val: %s err: %s", dataKey, string(rr.Value), err)
		} else {
			pathsData, err := allDomainPathsData(getDomainPathsDataKey(domainD.Id))
			if err == nil {
				domainD.Paths = append(domainD.Paths, pathsData...)
			}
			domainsData = append(domainsData, domainD)
		}
	}
	return domainsData, nil
}

//获取域名下所有定义路径数据
func allDomainPathsData(dataKey string) ([]*domain.Path, error) {
	var pathsData []*domain.Path
	ctx, cancel := context.WithTimeout(context.Background(), base.ReadTimeout)
	rsp, err := base.Cli.Get(ctx, dataKey, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return pathsData, err
	}
	for _, rr := range rsp.Kvs {
		path := new(domain.Path)
		err := json.Unmarshal(rr.Value, path)
		if err == nil {
			pathsData = append(pathsData, path)
		} else {
			base.Sys().Warnf("路径数据json解析失败 key : %s val: %s err: %s", dataKey, string(rr.Value), err)
		}
	}
	return pathsData, nil
}

func GetDomainPathDataByKey(dataKey string) (*domain.Path, error) {
	path := new(domain.Path)
	ctx, cancel := context.WithTimeout(context.Background(), base.ReadTimeout)
	rsp, err := base.Cli.Get(ctx, dataKey)
	cancel()
	if err != nil {
		return path, err
	}
	if rsp.Count > 0 {
		err = json.Unmarshal(rsp.Kvs[0].Value, path)
		if err != nil {
			base.Sys().Warnf("路径数据json解析失败 key : %s val: %s err: %s", dataKey, string(rsp.Kvs[0].Value), err)
		}
	}
	return path, nil
}

func WatchDomainsDataChange(e chan *clientv3.Event) {
	for {
		rch := base.Cli.Watch(context.Background(), base.DomainsDataPrefix, clientv3.WithPrefix())
		for rsp := range rch {
			for _, ev := range rsp.Events {
				e <- ev
			}
		}
	}
}

func WatchPathsDataChange(e chan *clientv3.Event) {
	for {
		rch := base.Cli.Watch(context.Background(), base.DomainPathsDataPrefix, clientv3.WithPrefix())
		for rsp := range rch {
			for _, ev := range rsp.Events {
				e <- ev
			}
		}
	}
}
