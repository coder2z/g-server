package manager

import (
	"errors"
	"github.com/coder2m/g-saber/xconsole"
	"github.com/coder2m/component/xcfg"
	"github.com/coder2m/component/xcfg/datasource/apollo"
	"github.com/coder2m/component/xcfg/datasource/etcdv3"
	"github.com/coder2m/component/xcfg/datasource/file"
	"net/url"
)

var (
	//ErrConfigAddr not xcfg
	ErrConfigAddr = errors.New("no xcfg... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	registry             map[string]DataSourceCreatorFunc
	//DefaultScheme ..
	DefaultScheme = `file`
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() xcfg.DataSource

func init() {
	registry = make(map[string]DataSourceCreatorFunc)
	Register(file.Register())
	Register(etcdv3.Register())
	Register(apollo.Register())
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator DataSourceCreatorFunc) {
	registry[scheme] = creator
}

//NewDataSource ..
func NewDataSource(configAddr string) (xcfg.DataSource, error) {
	if configAddr == "" {
		return nil, ErrConfigAddr
	}
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		DefaultScheme = urlObj.Scheme
	}
	creatorFunc, exist := registry[DefaultScheme]
	if !exist {
		return nil, ErrInvalidDataSource
	}
	xconsole.Greenf("Get xcfg from:", configAddr)
	source := creatorFunc()
	if source == nil {
		return nil, ErrInvalidDataSource
	}
	return source, nil
}
