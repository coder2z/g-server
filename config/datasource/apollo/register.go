package apollo

import (
	"flag"
	"github.com/myxy99/component/config"
	"github.com/myxy99/component/config/datasource/manager"
	"github.com/philchia/agollo/v4"
	"net/url"
)

// DataSourceApollo defines apollo scheme
const DataSourceApollo = "apollo"

func init() {
	manager.Register(DataSourceApollo, func() config.DataSource {
		var (
			configAddr = flag.String("config", "", "")
		)
		if *configAddr == "" {
			return nil
		}
		// configAddr is a string in this format:
		// apollo://ip:port?appId=XXX&cluster=XXX&namespaceName=XXX&key=XXX&accesskeySecret=XXX&insecureSkipVerify=XXX&cacheDir=XXX
		urlObj, err := url.Parse(*configAddr)
		if err != nil {
			return nil
		}
		apolloConf := agollo.Conf{
			AppID:              urlObj.Query().Get("appId"),
			Cluster:            urlObj.Query().Get("cluster"),
			NameSpaceNames:     []string{urlObj.Query().Get("namespaceName")},
			MetaAddr:           urlObj.Host,
			InsecureSkipVerify: true,
			AccesskeySecret:    urlObj.Query().Get("accesskeySecret"),
			CacheDir:           ".",
		}
		if urlObj.Query().Get("insecureSkipVerify") == "false" {
			apolloConf.InsecureSkipVerify = false
		}
		if urlObj.Query().Get("cacheDir") != "" {
			apolloConf.CacheDir = urlObj.Query().Get("cacheDir")
		}
		return NewDataSource(&apolloConf, urlObj.Query().Get("namespaceName"), urlObj.Query().Get("key"))
	})
}
