package file

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xflag"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func Register() (string, func() xcfg.DataSource) {
	return DataSourceFile, func() xcfg.DataSource {
		var (
			configAddr  = xflag.String("xcfg")
			watchConfig = xflag.Bool("watch")
		)
		if configAddr == "" {
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
}
