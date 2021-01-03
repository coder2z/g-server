/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 18:20
 **/
package xapp

import (
	"github.com/myxy99/component/pkg/xcast"
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
	debug           = true
)

func init() {
	var err error

	appName = xcfg.GetString("app.name")
	if appName == "" {
		appName = dAppName
	}

	if data := xcfg.GetString("app.debug"); data == "false" {
		debug = false
	}
	_ = os.Setenv("app.debug", xcast.ToString(debug))

	buildAppVersion = xcfg.GetString("app.version")
	if buildAppVersion == "" {
		buildAppVersion = dAppVersion
	}

	hostName, err = os.Hostname()
	if err != nil {
		hostName = "unknown"
	}

	startTime = time.Now().Format("2006-01-02 15:04:05")
	buildHost, _ = xnet.GetLocalIP()
	goVersion = runtime.Version()
}

// Name gets application name.
func Name() string {
	return appName
}

// Debug gets application debug.
func Debug() bool {
	return debug
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
