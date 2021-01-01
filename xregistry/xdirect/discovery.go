/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/1 22:44
 */
package xdirect

import (
	"fmt"
	"github.com/myxy99/component/xregistry"
	"strings"
)

type directDiscovery struct{}

func NewDiscovery() xregistry.Discovery {
	return &directDiscovery{}
}

// target 格式： "127.0.0.1:8000,127.0.0.1:8001"
func (d *directDiscovery) Discover(target string) (<-chan []xregistry.Instance, error) {
	endpoints := strings.Split(target, ",")
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoint")
	}
	var i []xregistry.Instance
	for _, addr := range endpoints {
		ins := xregistry.Instance{Address: addr}
		i = append(i, ins)
	}

	ch := make(chan []xregistry.Instance)
	go func() {
		ch <- i
	}()
	return ch, nil
}

func (d *directDiscovery) Close() {}
