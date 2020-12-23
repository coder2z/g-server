package file

import (
	"flag"
	"github.com/myxy99/component/config"
	"github.com/myxy99/component/pkg/xflag"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func Register() (string, func() config.DataSource) {
	return DataSourceFile, func() config.DataSource {
		var (
			configAddr  = xflag.String("config")
			watchConfig = xflag.Bool("watch")
		)
		flag.Parse()
		if configAddr == "" {
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
}
