# Component

```
    xcfg              这里主要是借鉴douyu框架得config解析
    xgoven            系统监控api
    xinvoker        
        -email        email
        -gorm         gorm
        -mongo        mongo
        -oss          云对象存储
        -redis        redis
    xlog              日志
    xmoniter          普罗米修斯
    xregistry         服务注册发现
    xtrace            链路追踪
    xgrpc             grpc 拦截器
    pkg
        -xcast        spf13/cast不足的函数封装
        -xcolor       颜色输出
        -xdefer       统一defer函数
        -xerrgroup    errgroup
        -xfile        文件操作封装
        -xflag        flag封装
        -xgp          协程池子
        -xmap         map操作封装
        -xnet         本机网络操作封装
        -xreminder    定时任务
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
	"github.com/coder2z/g-saber/xflag"
	"github.com/coder2z/component/xcfg"
	"github.com/coder2z/component/xcfg/datasource/manager"
	"github.com/coder2z/component/xinvoker"
	xgorm "github.com/coder2z/component/xinvoker/gorm"
	xredis "github.com/coder2z/component/xinvoker/redis"
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
	err = xcfg.LoadFromDataSource(data, toml.Unmarshal)
	if err != nil {
		panic(err)
	}

	xinvoker.Register(
		xgorm.Register("db"),
		xredis.Register("redis"))
	_ = xinvoker.Init()

	db := xgorm.Invoker("dev")
	rc := xredis.Invoker("dev")
	d, _ := db.DB()
	fmt.Println(d.Ping(), rc.Ping(context.Background()))
}


```

run

```bash
 go run main.go run -c=test.toml
```