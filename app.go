/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 18:20
 **/
package xapp

import (
	"fmt"
	"github.com/myxy99/component/pkg/xcolor"
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
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppName"), xcolor.Blue(Name()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("HostName"), xcolor.Blue(HostName()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppVersion"), xcolor.Blue(AppVersion()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildHost"), xcolor.Blue(BuildHost()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("StartTime"), xcolor.Blue(StartTime()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("GoVersion"), xcolor.Blue(GoVersion()))
}
