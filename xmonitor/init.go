/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:29
 **/
package xmonitor

import (
	cfg "github.com/myxy99/component/xcfg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"runtime"
	"time"
)

func init() {
	BuildInfoGauge.WithLabelValues(
		cfg.GetString("app.name"),
		cfg.GetString("app.mode"),
		cfg.GetString("app.version"),
		runtime.Version(),
		time.Now().String(),
	).Set(float64(time.Now().UnixNano() / 1e6))

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
}
