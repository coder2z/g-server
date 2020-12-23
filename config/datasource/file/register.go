package file

import (
	"flag"
	"github.com/myxy99/component/config"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func Register() (string, func() config.DataSource) {
	return DataSourceFile, func() config.DataSource {
		var (
			watchConfig = flag.Bool("watch", false, "")
			configAddr  = flag.String("config", "", "")
		)
		flag.Parse()
		if *configAddr == "" {
			return nil
		}
		return NewDataSource(*configAddr, *watchConfig)
	}
}
