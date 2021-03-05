/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 17:09
 **/
package xmongo

import (
	"github.com/coder2m/component/xcfg"
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
		Addrs:    []string{o.URL},
		Source:   o.Source,
		Username: o.User,
		Password: o.Password,
	}
	resp, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	mgo.SetDebug(o.Debug)
	// Optional. Switch the session to a monotonic behavior.
	resp.SetMode(mgo.Monotonic, true)
	return &client{resp}
}
