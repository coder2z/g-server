# g-server

Self-use microservice framework

![](https://img.shields.io/badge/windowns10-Development-d0d1d4)
![](https://img.shields.io/badge/golang-1.16-blue)
[![](https://img.shields.io/badge/godoc-reference-3C57C4)](https://pkg.go.dev/github.com/coder2z/g-server)
![](https://img.shields.io/badge/version-1.0.5-r)

## :rocket:Installation

`
go get -u github.com/coder2z/g-server
`

## :bell:Features

Encapsulate many tool functions and methods, and continue to update. Most of the methods are open source and other open
source projects.

```
xapp            = > runtime app info
datasource      = > xcfg data 
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

```env
set SERVER_APP_MODE     dev
set SERVER_APP_ID       79@sh$x3@
set SERVER_APP_DEBUG    false
```

```go
package main

import (
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xflag"
	"github.com/coder2z/g-server/xapp"
)

func main() {
	xflag.NewRootCommand(
		&xflag.CommandNode{
			Name: "main",
			Command: &xflag.Command{
				Use:                "main",
				DisableSuggestions: false,
				Run: func(cmd *xflag.Command, args []string) {
					RunApp()
				},
			},
			Flags: func(c *xflag.Command) {
				c.Flags().StringP("xcfg", "c", "", "配置文件")
				_ = c.MarkFlagRequired("config")
			},
		},
	)
	_ = xflag.Parse()
}

func RunApp() {
	xapp.PrintVersion()
	xconsole.Green("running")
}
```

### run

```bash
go build -o mainApp -ldflags "-X github.com/coder2z/g-server/xapp.appName=app_name -X github.com/coder2z/g-server/xapp.buildAppVersion=v1.0 -X github.com/coder2z/g-server/xapp.buildHost=`hostname`" main.go
./mainApp -c=test.toml
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
