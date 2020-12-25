/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:45
 **/
package xmonitor

import "github.com/prometheus/client_golang/prometheus"

var _registry = &Registry{
	Registerer: prometheus.DefaultRegisterer,
	namespace:  DefaultNamespace,
}

// NewCounter ...
func NewCounter(name string) prometheus.Counter {
	return _registry.NewCounter(name)
}

// NewCounterVec 计数器
func NewCounterVec(name string, labels []string) *prometheus.CounterVec {
	return _registry.NewCounterVec(name, labels)
}

// NewGauge ...
func NewGauge(name string) prometheus.Gauge {
	return _registry.NewGauge(name)
}

// NewGaugeVec gauge
func NewGaugeVec(name string, labels []string) *prometheus.GaugeVec {
	return _registry.NewGaugeVec(name, labels)
}

// NewHistogram ...
func NewHistogram(key string) prometheus.Histogram {
	return _registry.NewHistogram(key)
}

// NewHistogramVec ...
func NewHistogramVec(key string, labels []string) *prometheus.HistogramVec {
	return _registry.NewHistogramVec(key, labels)
}

// NewSummary ...
func NewSummary(name string) prometheus.Summary {
	return _registry.NewSummary(name)
}

// NewSummaryVec ...
func NewSummaryVec(name string, labels []string) *prometheus.SummaryVec {
	return _registry.NewSummaryVec(name, labels)
}
