package xapp

import (
	"fmt"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xcolor"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/g-saber/xtime"
	"github.com/coder2z/g-server/xversion"
	"os"
	"runtime"
)

var (
	DAppName  = "unknown app name"
	DHostName = "unknown host name"
	DHostIp   = "0.0.0.0"

	logo = ` 
   ____             ______ ______________  __ ___________ 
  / ___\   ______  /  ___// __ \_  __ \  \/ // __ \_  __ \
 / /_/  > /_____/  \___ \\  ___/|  | \/\   /\  ___/|  | \/
 \___  /          /____  >\___  >__|    \_/  \___  >__|   
/_____/                \/     \/                 \/ 

Hello, starting application ...
`
)

var info = appInfo{
	startTime:       xtime.Now().Format("2006-01-02 15:04:05"),
	goVersion:       runtime.Version(),
	hostName:        DHostName,
	hostIp:          DHostIp,
	AppName:         DAppName,
	BuildAppVersion: "v1.0.0",
	AppMode:         "dev",
	AppID:           "xxxxxx",
	Debug:           true,
}

// RegisterAppInfoCfg 注册app信息来自配置文件
func RegisterAppInfoCfg(key string) {
	info.hostName, _ = os.Hostname()
	info.hostIp, _ = xnet.GetLocalIP()
	info = *xcfg.UnmarshalWithExpect(key, &info).(*appInfo)
}

func init() {
	fmt.Println(xcolor.Blue(logo))
}

type appInfo struct {
	startTime string
	goVersion string
	hostName  string
	hostIp    string

	// build -X
	AppName         string
	BuildAppVersion string
	AppMode         string
	AppID           string
	Debug           bool
}

// Name gets application name.
func Name() string {
	return info.AppName
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return info.BuildAppVersion
}

//AppMode get AppMode
func AppMode() string {
	return info.AppMode
}

//HostIP get HostIP
func HostIP() string {
	return info.hostIp
}

// HostName get host name
func HostName() string {
	return info.hostName
}

//StartTime get start time
func StartTime() string {
	return info.startTime
}

//GoVersion get go version
func GoVersion() string {
	return info.goVersion
}

func AppId() string {
	return info.AppID
}

func Debug() bool {
	return info.Debug
}

func PrintVersion() {
	xlog.Infow("ApplicationInfo", "AppName:", Name())
	xlog.Infow("ApplicationInfo", "AppId:", AppId())
	xlog.Infow("ApplicationInfo", "AppMode:", AppMode())
	xlog.Infow("ApplicationInfo", "HostName:", HostName())
	xlog.Infow("ApplicationInfo", "HostIp:", HostIP())
	xlog.Infow("ApplicationInfo", "Debug:", Debug())
	xlog.Infow("ApplicationInfo", "AppVersion:", AppVersion())
	xlog.Infow("ApplicationInfo", "StartTime:", StartTime())
	xlog.Infow("ApplicationInfo", "G-server Version:", xversion.Version)
	xlog.Infow("ApplicationInfo", "GoVersion:", GoVersion())
}
