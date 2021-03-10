# Component

Self-use microservice framework

![](https://img.shields.io/badge/windowns10-Development-d0d1d4)
![](https://img.shields.io/badge/golang-1.16-blue)
[![](https://img.shields.io/badge/godoc-reference-3C57C4)](https://pkg.go.dev/github.com/coder2z/component)
![](https://img.shields.io/badge/version-1.0.5-r)

## :rocket:Installation

`
go get -u github.com/coder2z/component
`

## :bell:Features

Encapsulate many tool functions and methods, and continue to update. Most of the methods are open source and other open
source projects.

```
xapp            = > runtime app info
xcfg            = > config
xcode           = > err coder encapsulation
xgovern         = > system monitoring
xgrpc           = > grpc encapsulation
xinvoker        = > invoker
xmoniter        = > prometheus
xregistry       = > service registration discovery
xtrace          = > trace
xversion        = > frame version
```

## :anchor:Usage

### config file

```toml
[db]
    [db.dev]
        password = "root"
        dbName = "ndisk"

[redis]
    [redis.dev]
```

### main.go usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coder2z/g-saber/xflag"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xcfg/datasource/manager"
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

### run

```bash
 go run main.go run -c=test.toml
```

## :tada:Contribute code

Open source projects are inseparable from everyone’s support. If you have a good idea, encountered some bugs and fixed
them, and corrected the errors in the document, please submit a Pull Request~

1. Fork this project to your own repo
2. Clone the project in the past, that is, the project in your warehouse, to your local
3. Modify the code
4. Push to your own library after commit
5. Initiate a PR (pull request) request and submit it to the `provide` branch
6. Waiting to merge

## :closed_book:License

Distributed under MIT License, please see license file within the code for more details.