package xgorm

import (
	"github.com/coder2m/component/xcfg"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (i *dbInvoker) newDatabaseClient(o *options) (db *gorm.DB) {
	var err error
	db, err = gorm.Open(o.getDSN(), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   o.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	if o.Debug {
		db = db.Debug()
	}
	d, err := db.DB()
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(o.MaxOpenConnections)
	d.SetMaxIdleConns(o.MaxIdleConn)
	d.SetConnMaxLifetime(o.MaxConnectionLifeTime)
	//d.SetConnMaxIdleTime(o.MaxConnMaxIdleTime)
	return db
}

func (i *dbInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newDatabaseOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}
