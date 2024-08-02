package metrics_controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
)

type MemoryUsageGauge struct {
	Zone            string
	MemoryUsageDesc *prometheus.Desc
	MemoryUsage     float64
}

// Describe 向Prometheus描述收集的指标
func (g *MemoryUsageGauge) Describe(ch chan<- *prometheus.Desc) {
	ch <- g.MemoryUsageDesc
}

// Collect 收集指标数据并发送到Prometheus
func (g *MemoryUsageGauge) Collect(ch chan<- prometheus.Metric) {
	g.updateMemoryUsage()
	ch <- prometheus.MustNewConstMetric(
		g.MemoryUsageDesc,
		prometheus.GaugeValue,
		g.MemoryUsage,
	)
}

// NewMemoryUsageGauge 创建一个新的MemoryUsageGauge实例
func NewMemoryUsageGauge(zone string) *MemoryUsageGauge {

	return &MemoryUsageGauge{
		Zone: zone,
		MemoryUsageDesc: prometheus.NewDesc(
			"memory_usage_bytes",
			"系统内存使用量",
			nil,
			prometheus.Labels{"zone": zone},
		),
	}
}

// updateMemoryUsage 是一个用于更新内存使用量的方法。
// 它通过读取运行时的内存统计信息来获取当前分配的内存量，并将其更新到MemoryUsageGauge的实例中。
func (g *MemoryUsageGauge) updateMemoryUsage() {
	// 获取当前内存使用状况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// 更新内存使用量
	g.MemoryUsage = float64(m.Alloc)
}
