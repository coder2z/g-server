package file

import (
	"flag"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/pkg/xflag"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func Register() (string, func() xcfg.DataSource) {
	return DataSourceFile, func() xcfg.DataSource {
		var (
			configAddr  = xflag.String("xcfg")
			watchConfig = xflag.Bool("watch")
		)
		flag.Parse()
		if configAddr == "" {
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
}
