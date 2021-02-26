/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 17:35
 **/
package xgovern

import (
	"fmt"
	"github.com/coder2m/component/pkg/xnet"
)

type Config struct {
	Host    string `mapStructure:"host"`
	Port    int    `mapStructure:"port"`
	Network string `mapStructure:"network"`
}

type Option func(c *Config)

func DefaultConfig() *Config {
	host, port, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}
	return &Config{
		Host:    host,
		Port:    port,
		Network: "tcp",
	}
}

func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

func WithNetwork(network string) Option {
	return func(c *Config) {
		c.Network = network
	}
}

func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}
