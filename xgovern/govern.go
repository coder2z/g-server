/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 17:33
 **/
package xgovern

import (
	"context"
	"errors"
	"fmt"
	xapp "github.com/coder2m/component"
	"github.com/coder2m/component/xcfg"
	"github.com/coder2m/component/xcode"
	"github.com/coder2m/component/xlog"
	"github.com/coder2m/g-saber/xconsole"
	"github.com/coder2m/g-saber/xdefer"
	"github.com/coder2m/g-saber/xjson"
	"github.com/coder2m/g-saber/xnet"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"
	"time"
)

type healthStats struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Time     string `json:"time"`
	Err      string `json:"err"`
	Status   string `json:"status"`
}

type h map[string]func(w http.ResponseWriter, r *http.Request)

var (
	handle       *http.ServeMux
	server       *http.Server
	HandleFuncs  = make(h)
	governConfig *Config
	once         = sync.Once{}
)

func (hm h) Run(hs *http.ServeMux) {
	for s, f := range hm {
		hs.HandleFunc(s, f)
	}
}

func HandleFunc(p string, h func(w http.ResponseWriter, r *http.Request)) {
	HandleFuncs[p] = h
}

func init() {
	xcfg.OnChange(func(*xcfg.Configuration) {
		GovernReload()
	})
}

func init() {
	HandleFunc("/debug/pprof", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	HandleFunc("/debug/code/business", xcode.XCodeBusinessCodeHttp)
	HandleFunc("/debug/code/system", xcode.XCodeSystemCodeHttp)

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = xjson.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/debug/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		list := make([]string, 0)
		for s, _ := range HandleFuncs {
			list = append(list, s)
		}
		_ = xjson.NewEncoder(w).Encode(list)
	})

	HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		mm := xcfg.Traverse(".")
		w.WriteHeader(200)
		_ = xjson.NewEncoder(w).Encode(mm)
	})

	HandleFunc("/debug/health", func(w http.ResponseWriter, r *http.Request) {
		ip, _ := xnet.GetLocalIP()
		serverStats := healthStats{
			IP:       ip,
			Hostname: xapp.HostName(),
			Time:     time.Now().Format("2006-01-02 15:04:05"),
			Status:   "SUCCESS",
		}
		w.WriteHeader(200)
		_ = xjson.NewEncoder(w).Encode(serverStats)
	})
}

func GetServer() *http.ServeMux {
	if handle == nil {
		handle = http.NewServeMux()
	}
	return handle
}

func Run(opts ...Option) {
	once.Do(func() {
		c := GovernConfig()

		for _, opt := range opts {
			opt(c)
		}

		HandleFuncs.Run(GetServer())

		server = &http.Server{
			Addr:    c.Address(),
			Handler: handle,
		}

		xconsole.Greenf("govern serve init", fmt.Sprintf("%v/debug/list", c.Address()))

		xdefer.Register(func() error {
			return Shutdown()
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			xlog.Error("govern serve", xlog.FieldErr(err), xlog.FieldAddr(c.Address()))
		}
	})
}

func Shutdown() error {
	if server == nil {
		return errors.New("shutdown govern server")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		xlog.Error("shutdown govern server", xlog.FieldErr(err))
		return err
	}
	xconsole.Red("govern server shutdown")
	return nil
}

func GovernConfig() *Config {
	if governConfig == nil {
		governConfig = xcfg.UnmarshalWithExpect("app.govern", DefaultConfig()).(*Config)
	}
	return governConfig
}

func GovernReload(opts ...Option) {
	xconsole.Green("govern serve reload")
	_ = Shutdown()
	once = sync.Once{}
	governConfig = nil
	server = nil
	handle = nil
	Run(opts...)
}
