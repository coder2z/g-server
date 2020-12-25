/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:34
 **/
package xmonitor

import "github.com/prometheus/client_golang/prometheus"

type Registry struct {
	prometheus.Registerer
	namespace string
	subsystem string
}

// NewCounter ...
func (reg *Registry) NewCounter(name string) prometheus.Counter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
	})

	reg.MustRegister(counter)

	return counter
}

// NewCounterVec 计数器
func (reg *Registry) NewCounterVec(name string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
	}, labels)

	reg.MustRegister(counter)
	return counter
}

// NewGauge ...
func (reg *Registry) NewGauge(name string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
	})

	reg.MustRegister(gauge)
	return gauge
}

// NewGaugeVec gauge
func (reg *Registry) NewGaugeVec(name string, labels []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
	}, labels)
	reg.MustRegister(gauge)
	return gauge
}

// NewHistogram ...
func (reg *Registry) NewHistogram(key string) prometheus.Histogram {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        key,
		ConstLabels: nil,
		Buckets:     nil,
	})
	reg.MustRegister(histogram)
	return histogram
}

// NewHistogramVec ...
func (reg *Registry) NewHistogramVec(key string, labels []string) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        key,
		Help:        "summary of " + key,
		ConstLabels: nil,
		Buckets:     nil,
	}, labels)

	reg.MustRegister(histogram)
	return histogram
}

// NewSummary ...
func (reg *Registry) NewSummary(name string) prometheus.Summary {
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
		Objectives:  nil,
		MaxAge:      0,
		AgeBuckets:  0,
		BufCap:      0,
	})
	reg.MustRegister(summary)
	return summary
}

// NewSummaryVec ...
func (reg *Registry) NewSummaryVec(name string, labels []string) *prometheus.SummaryVec {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
		Objectives:  nil,
		MaxAge:      0,
		AgeBuckets:  0,
		BufCap:      0,
	}, labels)
	reg.MustRegister(summary)
	return summary
}

type TimerVec struct {
	*prometheus.HistogramVec
	*prometheus.CounterVec
}

func (reg *Registry) NewTimer(name string, labels []string) *TimerVec {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
		Buckets:     nil,
	}, labels)

	reg.MustRegister(histogram)

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   reg.namespace,
		Subsystem:   reg.subsystem,
		Name:        name,
		ConstLabels: nil,
	}, labels)

	reg.MustRegister(counter)

	return &TimerVec{
		HistogramVec: histogram,
		CounterVec:   counter,
	}
}
