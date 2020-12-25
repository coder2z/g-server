/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:48
 **/
package xmonitor

var (
	// TypeHTTP ...
	TypeHTTP = "http"
	// TypeGRPCUnary ...
	TypeGRPCUnary = "unary"
	// TypeGRPCStream ...
	TypeGRPCStream = "stream"
	// TypeRedis ...
	TypeRedis = "redis"
	// TypeGorm ...
	TypeGorm = "gorm"
	// TypeRocketMQ ...
	TypeRocketMQ = "rocketmq"
	// TypeWebsocket ...
	TypeWebsocket = "ws"
	// TypeMySQL ...
	TypeMySQL = "mysql"
	// CodeJob
	CodeJobSuccess = "ok"
	// CodeJobFail ...
	CodeJobFail = "fail"
	// CodeJobReentry ...
	CodeJobReentry = "reentry"
	// CodeCache
	CodeCacheMiss = "miss"
	// CodeCacheHit ...
	CodeCacheHit = "hit"
	// Namespace
	DefaultNamespace = "xmonitor"
)

var (
	// ServerHandleCounter ...
	ServerHandleCounter = NewCounterVec("server_handle_total", []string{"type", "method", "peer", "code"})

	// ServerHandleHistogram ...
	ServerHandleHistogram = NewHistogramVec("server_handle_seconds", []string{"type", "method", "peer"})

	// ClientHandleCounter ...
	ClientHandleCounter = NewCounterVec("client_handle_total", []string{"type", "name", "method", "peer", "code"})

	// ClientHandleHistogram ...
	ClientHandleHistogram = NewHistogramVec("client_handle_seconds", []string{"type", "name", "method", "peer"})

	// JobHandleCounter ...
	JobHandleCounter = NewCounterVec("job_handle_total", []string{"type", "name", "code"})

	// JobHandleHistogram ...
	JobHandleHistogram = NewHistogramVec("job_handle_seconds", []string{"type", "name"})

	LibHandleHistogram = NewHistogramVec("lib_handle_seconds", []string{"type", "method", "address"})

	// LibHandleCounter ...
	LibHandleCounter = NewCounterVec("lib_handle_total", []string{"type", "method", "address", "code"})

	LibHandleSummary = NewSummaryVec("lib_handle_stats", []string{"name", "status"})

	// CacheHandleCounter ...
	CacheHandleCounter = NewCounterVec("cache_handle_total", []string{"type", "name", "action", "code"})

	// CacheHandleHistogram ...
	CacheHandleHistogram = NewHistogramVec("cache_handle_seconds", []string{"type", "name", "action"})

	// BuildInfoGauge ...
	BuildInfoGauge = NewGaugeVec("build_info", []string{"name", "mode", "app_version", "go_version", "start_time"})
)
