/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package xmongo

import (
	"time"
)

type options struct {
	Addrs          []string      `mapStructure:"addrs"`
	Timeout        time.Duration `mapStructure:"timeout"`
	Database       string        `mapStructure:"database"`
	ReplicaSetName string        `mapStructure:"replica_set_name"`
	Source         string        `mapStructure:"source"`
	Service        string        `mapStructure:"service"`
	ServiceHost    string        `mapStructure:"service_host"`
	Mechanism      string        `mapStructure:"mechanism"`
	Username       string        `mapStructure:"username"`
	Password       string        `mapStructure:"password"`
	PoolLimit      int           `mapStructure:"pool_limit"`
	PoolTimeout    time.Duration `mapStructure:"pool_timeout"`
	ReadTimeout    time.Duration `mapStructure:"read_timeout"`
	WriteTimeout   time.Duration `mapStructure:"write_timeout"`
	AppName        string        `mapStructure:"app_name"`
	FailFast       bool          `mapStructure:"fail_fast"`
	Direct         bool          `mapStructure:"direct"`
	MinPoolSize    int           `mapStructure:"min_pool_size"`
	MaxIdleTimeMS  int           `mapStructure:"max_idle_time_ms"`
	Debug          bool          `mapStructure:"debug"`
	URL            []string      `mapStructure:"url"`
	User           string        `mapStructure:"user"`
}

func newMongoOptions() *options {
	return &options{
		Addrs:          []string{},
		Timeout:        0,
		Database:       "",
		ReplicaSetName: "",
		Source:         "",
		Service:        "",
		ServiceHost:    "",
		Mechanism:      "",
		Username:       "",
		Password:       "",
		PoolLimit:      0,
		PoolTimeout:    0,
		ReadTimeout:    0,
		WriteTimeout:   0,
		AppName:        "",
		FailFast:       false,
		Direct:         false,
		MinPoolSize:    0,
		MaxIdleTimeMS:  0,
		Debug:          false,
		URL:            []string{},
		User:           "",
	}
}
