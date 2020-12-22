package manager

import (
	"errors"
	"github.com/myxy99/component/config"
	"net/url"
)

var (
	//ErrConfigAddr not config
	ErrConfigAddr = errors.New("no config... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	registry             map[string]DataSourceCreatorFunc
	//DefaultScheme ..
	DefaultScheme string
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() config.DataSource

func init() {
	registry = make(map[string]DataSourceCreatorFunc)
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator DataSourceCreatorFunc) {
	registry[scheme] = creator
}

// CreateDataSource creates a dataSource witch has been registered
// func CreateDataSource(scheme string) (conf.DataSource, error) {
// 	creatorFunc, exist := registry[scheme]
// 	if !exist {
// 		return nil, ErrInvalidDataSource
// 	}
// 	return creatorFunc(), nil
// }

//NewDataSource ..
func NewDataSource(configAddr string) (config.DataSource, error) {
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
	return creatorFunc(), nil
}
