/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package xgorm

import (
	"fmt"
	"github.com/coder2m/component/xinvoker"
	"gorm.io/gorm"
	"sync"
)

var db *dbInvoker

func Register(k string) xinvoker.Invoker {
	db = &dbInvoker{key: k}
	return db
}

func Invoker(key string) *gorm.DB {
	if val, ok := db.instances.Load(key); ok {
		return val.(*gorm.DB)
	}
	panic(fmt.Sprintf("no db(%s) invoker found", key))
}

type dbInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *dbInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		db := i.newDatabaseClient(cfg)
		i.instances.Store(name, db)
	}
	return nil
}

func (i *dbInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		db := i.newDatabaseClient(cfg)
		i.instances.Store(name, db)
	}
	return nil
}

func (i *dbInvoker) Close(opts ...xinvoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		db, _ := value.(*gorm.DB).DB()
		_ = db.Close()
		i.instances.Delete(key)
		return true
	})
	return nil
}
