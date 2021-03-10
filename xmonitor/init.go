/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:29
 **/
package xmonitor

import (
	"github.com/coder2z/component/xapp"
	cfg "github.com/coder2z/g-saber/xcfg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func init() {
	BuildInfoGauge.WithLabelValues(
		xapp.Name(),
		cfg.GetString("app.mode"),
		xapp.AppVersion(),
		xapp.GoVersion(),
		xapp.StartTime(),
	).Set(float64(time.Now().UnixNano() / 1e6))
}

func MonitorPrometheusHttp(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
