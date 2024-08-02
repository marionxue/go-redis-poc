package metrics_controller

import "github.com/prometheus/client_golang/prometheus"

// CPULoadSummary 结构体，用于管理CPU负载的摘要指标
type CPULoadSummary struct {
	Summary prometheus.Summary
}

// NewCPULoadSummary 创建一个新的CPULoadSummary实例
func NewCPULoadSummary(zone string) *CPULoadSummary {
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "cpu_load_summary",
		Help:       "CPU负载摘要",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // 定义摘要的分位数和允许的误差
		ConstLabels: prometheus.Labels{
			"zone": zone,
		},
	})

	return &CPULoadSummary{
		Summary: summary,
	}
}

// Describe 实现prometheus.Collector接口中的Describe方法
func (s *CPULoadSummary) Describe(ch chan<- *prometheus.Desc) {
	s.Summary.Describe(ch)
}

// Collect 实现prometheus.Collector接口中的Collect方法
func (s *CPULoadSummary) Collect(ch chan<- prometheus.Metric) {
	s.Summary.Collect(ch)
}

// RecordCPULoad 记录CPU负载
func (s *CPULoadSummary) RecordCPULoad(load float64) {
	s.Summary.Observe(load)
}
