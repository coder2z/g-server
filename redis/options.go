/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package redis

import (
	"time"
)

type options struct {
	Network            string        `yaml:"network"`
	Addr               string        `yaml:"addr"`
	Username           string        `yaml:"username"`
	Password           string        `yaml:"password"`
	DB                 int           `yaml:"db"`
	MaxRetries         int           `yaml:"maxRetries"`
	MinRetryBackoff    time.Duration `yaml:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `yaml:"maxRetryBackoff"`
	DialTimeout        time.Duration `yaml:"dialTimeout"`
	ReadTimeout        time.Duration `yaml:"readTimeout"`
	WriteTimeout       time.Duration `yaml:"writeTimeout"`
	PoolSize           int           `yaml:"poolSize"`
	MinIdleConns       int           `yaml:"minIdleConn"`
	MaxConnAge         time.Duration `yaml:"maxConnAge"`
	PoolTimeout        time.Duration `yaml:"poolTimeout"`
	IdleTimeout        time.Duration `yaml:"idleTimeout"`
	IdleCheckFrequency time.Duration `yaml:"idleCheckFrequency"`
	readOnly           bool          `yaml:"readOnly"`
}

func newRedisOptions() *options {
	return &options{
		Addr:     "127.0.0.1:6379",
		Username: "",
		Password: "",
		DB:       0,
	}
}
