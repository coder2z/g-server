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
)

func init() {
	if appName == xcfg.GetString("app.name") {
		if appName == "" {
			appName = dAppName
		}
	}

	if buildAppVersion == xcfg.GetString("app.version") {
		if buildAppVersion == "" {
			buildAppVersion = dAppVersion
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = time.Now().Format("2006-01-02 15:04:05")
	buildHost, _ = xnet.GetLocalIP()
	goVersion = runtime.Version()
}

// Name gets application name.
func Name() string {
	return appName
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

//BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

// HostName get host name
func HostName() string {
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

func PrintVersion() {
	xconsole.Blue(fmt.Sprintf("%-40v", "——————————————————"))
	xconsole.Greenf("AppName", Name())
	xconsole.Greenf("HostName", HostName())
	xconsole.Greenf("AppVersion", AppVersion())
	xconsole.Greenf("BuildHost", BuildHost())
	xconsole.Greenf("StartTime", StartTime())
	xconsole.Greenf("GoVersion", GoVersion())
	xconsole.Blue(fmt.Sprintf("%-40v", "——————————————————"))
}
