package xapp

import (
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/g-saber/xtime"
	"github.com/coder2z/g-server/xversion"
	"os"
	"runtime"
)

const (
	dAppName = "unknown app name"
	dHostIp  = "0.0.0.0"
)

func init() {
	xconsole.Blue(`   _____ ____  _____  ______ _____  ___  ______`)
	xconsole.Blue(`  / ____/ __ \|  __ \|  ____|  __ \|__ \|___  /`)
	xconsole.Blue(` | |   | |  | | |  | | |__  | |__) |  ) |  / /`)
	xconsole.Blue(` | |   | |  | | |  | |  __| |  _  /  / /  / / `)
	xconsole.Blue(` | |___| |__| | |__| | |____| | \ \ / /_ / /__ `)
	xconsole.Blue(`  \_____\____/|_____/|______|_|  \_|____/_____|`)
	xconsole.Blue(`									--version = ` + xversion.Version)
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
	xconsole.Greenf("app name:", Name())
	xconsole.Greenf("app id:", AppId())
	xconsole.Greenf("app mode:", AppMode())
	xconsole.Greenf("host name:", HostName())
	xconsole.Greenf("host ip:", HostIP())
	xconsole.Greenf("debug:", Debug())
	xconsole.Greenf("app version:", AppVersion())
	xconsole.Greenf("build host:", BuildHost())
	xconsole.Greenf("start time:", StartTime())
	xconsole.Greenf("go version:", GoVersion())
}
