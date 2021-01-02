package etcdv3

import (
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/xcfg"
	"net/url"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// DataSourceEtcd defines etcd scheme
const DataSourceEtcd = "etcd"

func Register() (string, func() xcfg.DataSource) {
	return DataSourceEtcd, func() xcfg.DataSource {
		var (
			configAddr = xflag.String("xcfg")
		)
		if configAddr == "" {
			return nil
		}
		// configAddr is a string in this format:
		// etcd://ip:port?username=XXX&password=XXX&key=key

		urlObj, err := url.Parse(configAddr)
		if err != nil {
			return nil
		}
		etcdConf := clientv3.Config{
			DialKeepAliveTime:    10 * time.Second,
			DialKeepAliveTimeout: 3 * time.Second,
		}
		etcdConf.Endpoints = []string{urlObj.Host}
		etcdConf.Username = urlObj.Query().Get("username")
		etcdConf.Password = urlObj.Query().Get("password")
		client, err := clientv3.New(etcdConf)
		if err != nil {
			return nil
		}
		return NewDataSource(client, urlObj.Query().Get("key"))
	}
}
