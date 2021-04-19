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
	// ServerHandleCounter ...	指标: 服务类型，调用方法，客户端标识，返回的状态码
	ServerHandleCounter = NewCounterVec("server_handle_total", []string{"type", "name", "method", "peer", "code"})

	// ServerErrorCounter ...	指标: 服务类型，调用方法，客户端标识，返回的状态码
	ServerErrorCounter = NewCounterVec("server_error_total", []string{"type", "name", "method", "peer", "code"})

	// ServerHandleHistogram ...
	ServerHandleHistogram = NewHistogramVec("server_handle_seconds", []string{"type", "name", "method", "peer"})

	// ClientHandleCounter ... 	指标: 客户端类型，客户端名称，调用方法，目标，返回的状态码
	ClientHandleCounter = NewCounterVec("client_handle_total", []string{"type", "name", "method", "peer", "code"})

	// ClientHandleHistogram ...
	ClientHandleHistogram = NewHistogramVec("client_handle_seconds", []string{"type", "name", "method", "peer"})

	// JobHandleCounter ...	指标: 类型，任务名，执行状态码
	JobHandleCounter = NewCounterVec("job_handle_total", []string{"type", "name", "code"})

	// JobHandleHistogram ...
	JobHandleHistogram = NewHistogramVec("job_handle_seconds", []string{"type", "name"})

	// LibHandleHistogram ...	 指标: 类型，指令，address
	LibHandleHistogram = NewHistogramVec("lib_handle_seconds", []string{"type", "name", "method", "address"})

	//	LibHandleCounter ...
	LibHandleCounter = NewCounterVec("lib_handle_total", []string{"type", "name", "method", "address", "code"})

	//	LibHandleSummary
	LibHandleSummary = NewSummaryVec("lib_handle_stats", []string{"name", "status"})

	// CacheHandleCounter ...	指标: 类型，缓存名
	CacheHandleCounter = NewCounterVec("cache_handle_total", []string{"type", "name", "action", "code"})

	// CacheHandleHistogram ...
	CacheHandleHistogram = NewHistogramVec("cache_handle_seconds", []string{"type", "name", "action"})

	// BuildInfoGauge ...	版本信息指标
	BuildInfoGauge = NewGaugeVec("build_info", []string{"name", "mode", "instance", "app_version", "go_version", "start_time"})
)
