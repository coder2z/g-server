package xmongo

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xlog"
	"github.com/globalsign/mgo"
)

type MongoImp interface {
	Connect(db, collection string) (*mgo.Session, *mgo.Collection)
	Insert(db, collection string, docs ...interface{}) error
	FindOne(db, collection string, query, selector, result interface{}) error
	FindAll(db, collection string, query, selector, result interface{}) error
	Update(db, collection string, query, update interface{}) error
	Remove(db, collection string, query interface{}) error
}

type client struct {
	m *mgo.Session
}

func (c *client) Connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := c.m.Copy()
	return s, s.DB(db).C(collection)
}

func (c *client) Insert(db, collection string, docs ...interface{}) error {
	ms, mc := c.Connect(db, collection)
	defer ms.Close()
	return mc.Insert(docs...)
}

func (c *client) FindOne(db, collection string, query, selector, result interface{}) error {
	ms, mc := c.Connect(db, collection)
	defer ms.Close()
	return mc.Find(query).Select(selector).One(result)
}

func (c *client) FindAll(db, collection string, query, selector, result interface{}) error {
	ms, mc := c.Connect(db, collection)
	defer ms.Close()
	return mc.Find(query).Select(selector).All(result)
}

func (c *client) Update(db, collection string, query, update interface{}) error {
	ms, mc := c.Connect(db, collection)
	defer ms.Close()
	return mc.Update(query, update)
}

func (c *client) Remove(db, collection string, query interface{}) error {
	ms, mc := c.Connect(db, collection)
	defer ms.Close()
	return mc.Remove(query)
}

func (i *mongoInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newMongoOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}

func (i *mongoInvoker) new(o *options) MongoImp {
	dialInfo := &mgo.DialInfo{
		Addrs:          o.Addrs,
		Timeout:        o.Timeout,
		Database:       o.Database,
		ReplicaSetName: o.ReplicaSetName,
		Source:         o.Source,
		Service:        o.Service,
		ServiceHost:    o.ServiceHost,
		Mechanism:      o.Mechanism,
		Username:       o.Username,
		Password:       o.Password,
		PoolLimit:      o.PoolLimit,
		PoolTimeout:    o.PoolTimeout,
		ReadTimeout:    o.ReadTimeout,
		WriteTimeout:   o.WriteTimeout,
		AppName:        o.AppName,
		FailFast:       o.FailFast,
		Direct:         o.Direct,
		MinPoolSize:    o.MinPoolSize,
		MaxIdleTimeMS:  o.MaxIdleTimeMS,
	}
	resp, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XInvoker"),
			xlog.FieldMethod("XInvoker.XMongo.new"),
			xlog.FieldDescription("new mongo error"),
			xlog.FieldErr(err),
		)
	}
	mgo.SetDebug(o.Debug)
	// Optional. Switch the session to a monotonic behavior.
	resp.SetMode(mgo.Monotonic, true)
	return &client{resp}
}
