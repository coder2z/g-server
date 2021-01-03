/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 18:20
 **/
package xapp

import (
	"fmt"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xnet"
	"github.com/myxy99/component/xcfg"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	dAppName    = "MyApp"
	dAppVersion = "v0.1.0"
)

var (
	startTime       string
	goVersion       string
	appName         string
	hostName        string
	buildAppVersion string
	buildHost       string
	debug           = true
	one             = sync.Once{}
)

// Name gets application name.
func Name() string {
	if appName == "" {
		if appName = xcfg.GetString("app.name"); appName == "" {
			appName = dAppName
		}
	}
	return appName
}

// Debug gets application debug.
func Debug() bool {
	if data := xcfg.GetString("app.debug"); data == "false" {
		debug = false
	}
	one.Do(func() {
		xconsole.ResetDebug(debug)
	})
	return debug
}

//AppVersion get buildAppVersion
func AppVersion() string {
	if buildAppVersion == "" {
		if buildAppVersion = xcfg.GetString("app.version"); buildAppVersion == "" {
			buildAppVersion = dAppVersion
		}
	}
	return buildAppVersion
}

//BuildHost get buildHost
func BuildHost() string {
	if buildHost == "" {
		var err error
		if buildHost, err = xnet.GetLocalIP(); err != nil {
			hostName = "0.0.0.0"
		}
	}
	return buildHost
}

// HostName get host name
func HostName() string {
	if hostName == "" {
		var err error
		if hostName, err = os.Hostname(); err != nil {
			hostName = "unknown"
		}
	}
	return hostName
}

//StartTime get start time
func StartTime() string {
	if startTime == "" {
		startTime = time.Now().Format("2006-01-02 15:04:05")
	}
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	if goVersion == "" {
		goVersion = runtime.Version()
	}
	return goVersion
}

func PrintVersion() {
	xconsole.Blue(fmt.Sprintf("%-40v", "——————————————————"))
	xconsole.Greenf("app name:", Name())
	xconsole.Greenf("host name:", HostName())
	xconsole.Greenf("app debug:", Debug())
	xconsole.Greenf("app version:", AppVersion())
	xconsole.Greenf("build host:", BuildHost())
	xconsole.Greenf("start time:", StartTime())
	xconsole.Greenf("go version:", GoVersion())
	xconsole.Blue(fmt.Sprintf("%-40v", "——————————————————"))
}
