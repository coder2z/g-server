/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:29
 **/
package xmonitor

import (
	"fmt"
	xapp "github.com/coder2z/component"
	"github.com/coder2z/g-saber/xconsole"
	cfg "github.com/coder2z/component/xcfg"
	"github.com/coder2z/component/xgovern"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func Run() {
	BuildInfoGauge.WithLabelValues(
		xapp.Name(),
		cfg.GetString("app.mode"),
		xapp.AppVersion(),
		xapp.GoVersion(),
		xapp.StartTime(),
	).Set(float64(time.Now().UnixNano() / 1e6))

	xgovern.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	xconsole.Greenf("prometheus monitor init", fmt.Sprintf("%v/metrics", xgovern.GovernConfig().Address()))
}
