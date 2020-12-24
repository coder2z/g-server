# Component

```
    config          这里主要是借鉴douyu框架得config解析
    email           邮件
    gorm            orm gorm v2
    log             日志
    mongo           mongoDB
    oss             云对象存储
    redis           redis
    pkg
      -xcast        spf13/cast不足的函数封装
      -xcolor       颜色输出
      -xfile        文件操作封装
      -xflag        flag封装
      -xgp          协程池子
      -xmap         map操作封装
      -xsignals     信号     
      -xtransform   数组转换封装   
      -xvalidator   验证
```

方便开发，自用

很多东西都是来自开源的项目。


## Example

配置文件

```toml
[db]
[db.dev]
password = "root"
dbName = "ndisk"

[redis]
[redis.dev]
```

使用

```go
/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/23 15:13
 **/
package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	invoker "github.com/myxy99/component"
	"github.com/myxy99/component/config"
	"github.com/myxy99/component/config/datasource/manager"
	database "github.com/myxy99/component/gorm"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/redis"
)

func main() {
	xflag.Register(
		xflag.CommandNode{
			Name: "run",
			Command: &xflag.Command{
				Use:   "run ",
				Short: "run your app",
				Run: func(cmd *xflag.Command, args []string) {
					RunApp()
				},
			},
			Flags: func(c *xflag.Command) {
				c.Flags().StringP("config", "c", "", "配置文件")
				_ = c.MarkFlagRequired("config")
			},
		},
	)

	_ = xflag.Parse()

}

func RunApp() {
	data, err := manager.NewDataSource(xflag.NString("run", "config"))
	if err != nil {
		panic(err)
	}
	err = config.LoadFromDataSource(data, toml.Unmarshal)
	if err != nil {
		panic(err)
	}

	invoker.Register(
		database.Register("db"),
		redis.Register("redis"))
	_ = invoker.Init()

	db := database.Invoker("dev")
	rc := redis.Invoker("dev")
	d, _ := db.DB()
	fmt.Println(d.Ping(), rc.Ping(context.Background()))
}

```

run

```bash
 go run main.go run -c=test.toml
```