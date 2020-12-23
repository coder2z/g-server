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
package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/myxy99/component"
	"github.com/myxy99/component/config"
	"github.com/myxy99/component/config/datasource/manager"
	database "github.com/myxy99/component/gorm"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/redis"
	"os"
)

func main() {
	_ = os.Setenv("CONFIG", "test.toml")
	xflag.Register(&xflag.StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  "CONFIG",
		Default: "",
	})

	xflag.Register(&xflag.BoolFlag{
		Name:    "watch",
		Usage:   "--watch, watch config change event",
		Default: false,
		EnvVar:  "CONFIG_WATCH",
	})

	_ = xflag.Parse()
    
    //获取配置数据源
	data, err := manager.NewDataSource(xflag.String("config"))
	if err != nil {
		panic(err)
	}

	//解析配置
	err = config.LoadFromDataSource(data, toml.Unmarshal)
	if err != nil {
		panic(err)
	}

	//需要初始化的东西
	invoker.Register(
		database.Register("db"),
		redis.Register("redis"))

	//初始化
	_ = invoker.Init()
    
    //获取
	db := database.Invoker("dev")
	rc := redis.Invoker("dev")

    //使用
	d, _ := db.DB()
	fmt.Println(d.Ping(), rc.Ping(context.Background()))

}

```