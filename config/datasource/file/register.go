package file

import (
	"flag"
	"github.com/myxy99/component/config"
	"github.com/myxy99/component/config/datasource/manager"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func init() {
	manager.Register(DataSourceFile, func() config.DataSource {
		var (
			watchConfig = flag.Bool("watch", false, "")
			configAddr  = flag.String("config", "", "")
		)
		if *configAddr == "" {
			return nil
		}
		return NewDataSource(*configAddr, *watchConfig)
	})
	manager.DefaultScheme = DataSourceFile
}
