package xapp

import (
	"fmt"
	"github.com/coder2z/g-saber/xcolor"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/g-saber/xtime"
	"github.com/coder2z/g-server/xversion"
	"os"
	"runtime"
)

const (
	dAppName = "unknown app name"
	dHostIp  = "0.0.0.0"

	logo = ` 
   ____             ______ ______________  __ ___________ 
  / ___\   ______  /  ___// __ \_  __ \  \/ // __ \_  __ \
 / /_/  > /_____/  \___ \\  ___/|  | \/\   /\  ___/|  | \/
 \___  /          /____  >\___  >__|    \_/  \___  >__|   
/_____/                \/     \/                 \/ 

Hello, starting application ...
`
)

func init() {
	fmt.Println(xcolor.Blue(logo))
	startTime = xtime.Now().Format("2006-01-02 15:04:05")
	goVersion = runtime.Version()

	name, err := os.Hostname()
	if err != nil {
		name = dAppName
	}
	hostName = name

	ip, err := xnet.GetLocalIP()
	if err != nil {
		ip = dHostIp
	}
	hostIp = ip

	appMode = os.Getenv(`SERVER_APP_MODE`)
	appId = os.Getenv(`SERVER_APP_ID`)
	envDebug := os.Getenv(`SERVER_APP_DEBUG`)
	if envDebug != "" {
		debug = envDebug
	}
}

var (
	startTime string
	goVersion string
	hostName  string
	hostIp    string

	// build -X
	appName         string
	buildAppVersion string
	buildHost       string

	// env
	appMode string
	appId   string
	debug   = "true"
)

// Name gets application name.
func Name() string {
	return appName
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

//AppMode get AppMode
func AppMode() string {
	return appMode
}

//HostIP get HostIP
func HostIP() string {
	return hostIp
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

func AppId() string {
	return appId
}

func Debug() bool {
	return debug == "true"
}

func PrintVersion() {
	xlog.Infow("ApplicationInfo", "AppName:", Name())
	xlog.Infow("ApplicationInfo", "AppId:", AppId())
	xlog.Infow("ApplicationInfo", "AppMode:", AppMode())
	xlog.Infow("ApplicationInfo", "HostName:", HostName())
	xlog.Infow("ApplicationInfo", "HostIp:", HostIP())
	xlog.Infow("ApplicationInfo", "Debug:", Debug())
	xlog.Infow("ApplicationInfo", "AppVersion:", AppVersion())
	xlog.Infow("ApplicationInfo", "BuildHost:", BuildHost())
	xlog.Infow("ApplicationInfo", "StartTime:", StartTime())
	xlog.Infow("ApplicationInfo", "G-server Version:", xversion.Version)
	xlog.Infow("ApplicationInfo", "GoVersion:", GoVersion())
}
