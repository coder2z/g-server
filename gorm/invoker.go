/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 13:52
 **/
package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	invoker "github.com/myxy99/component"
	"sync"
)

var db *dbInvoker

func Register(k string) invoker.Invoker {
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
	invoker.Base
	instances sync.Map
	key       string
}

func (i *dbInvoker) Init(opts ...invoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		db := i.newDatabaseClient(cfg)
		i.instances.Store(name, db)
	}
	return nil
}

func (i *dbInvoker) Reload(opts ...invoker.Option) error {
	for name, cfg := range i.loadConfig() {
		db := i.newDatabaseClient(cfg)
		i.instances.Store(name, db)
	}
	return nil
}

func (i *dbInvoker) Close(opts ...invoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		_ = value.(*gorm.DB).Close()
		i.instances.Delete(key)
		return true
	})
	return nil
}
