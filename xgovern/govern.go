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
	iJson "github.com/json-iterator/go"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xconsole"
	"github.com/myxy99/component/pkg/xdefer"
	"github.com/myxy99/component/pkg/xnet"
	"github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xlog"
	"net/http"
	"net/http/pprof"
	"os"
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
	handle      *http.ServeMux
	server      *http.Server
	HandleFuncs = make(h)
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
	HandleFunc("/debug/pprof", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = iJson.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/debug/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		list := make([]string, 0)
		for s, _ := range HandleFuncs {
			list = append(list, s)
		}
		_ = iJson.NewEncoder(w).Encode(list)
	})

	HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		mm := xcfg.Traverse(".")
		w.WriteHeader(200)
		_ = iJson.NewEncoder(w).Encode(mm)
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
		_ = iJson.NewEncoder(w).Encode(serverStats)
	})
}

func GetServer() *http.ServeMux {
	if handle == nil {
		handle = http.NewServeMux()
	}
	return handle
}

func Run(opts ...Option) {
	c := xcfg.UnmarshalWithExpect("app.govern", DefaultConfig()).(*Config)

	for _, opt := range opts {
		opt(c)
	}

	HandleFuncs.Run(GetServer())

	server = &http.Server{
		Addr:    c.Address(),
		Handler: handle,
	}

	xconsole.Greenf("govern serve init:", fmt.Sprintf("%v/debug/list", c.Address()))

	xdefer.Register(func() error {
		return Shutdown()
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		xlog.Errorw("govern serve", xlog.String("error", err.Error()), xlog.FieldAddr(c.Address()))
	}
}

func Shutdown() error {
	if server == nil {
		return errors.New("shutdown govern server")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		xlog.Errorw("shutdown govern server", xlog.FieldErr(err))
		return err
	}
	xconsole.Red("govern server shutdown ~")
	return nil
}
