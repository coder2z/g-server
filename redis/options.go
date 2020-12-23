/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package redis

import (
	"time"
)

type options struct {
	Network            string        `mapStructure:"network"`
	Addr               string        `mapStructure:"addr"`
	Username           string        `mapStructure:"username"`
	Password           string        `mapStructure:"password"`
	DB                 int           `mapStructure:"db"`
	MaxRetries         int           `mapStructure:"maxRetries"`
	MinRetryBackoff    time.Duration `mapStructure:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `mapStructure:"maxRetryBackoff"`
	DialTimeout        time.Duration `mapStructure:"dialTimeout"`
	ReadTimeout        time.Duration `mapStructure:"readTimeout"`
	WriteTimeout       time.Duration `mapStructure:"writeTimeout"`
	PoolSize           int           `mapStructure:"poolSize"`
	MinIdleConns       int           `mapStructure:"minIdleConn"`
	MaxConnAge         time.Duration `mapStructure:"maxConnAge"`
	PoolTimeout        time.Duration `mapStructure:"poolTimeout"`
	IdleTimeout        time.Duration `mapStructure:"idleTimeout"`
	IdleCheckFrequency time.Duration `mapStructure:"idleCheckFrequency"`
	readOnly           bool          `mapStructure:"readOnly"`
}

func newRedisOptions() *options {
	return &options{
		Addr:     "127.0.0.1:6379",
		Username: "",
		Password: "",
		DB:       0,
	}
}
