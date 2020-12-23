/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 12:02
 **/
package test

import (
	"fmt"
	"github.com/myxy99/component/config/datasource/manager"
	"testing"
)

func TestLocConfig(t *testing.T) {
	data, err := manager.NewDataSource("test.toml")
	fmt.Println(data, err)
}
