package metrics_controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync/atomic"
)

// APIRequestCounter 结构体，用于管理API请求次数的监控
type APIRequestCounter struct {
	Zone           string
	APIRequestDesc *prometheus.Desc
	requestCount   uint64
}

// NewAPIRequestCounter 创建一个新的APIRequestCounter实例
func NewAPIRequestCounter(zone string) *APIRequestCounter {
	return &APIRequestCounter{
		Zone: zone,
		APIRequestDesc: prometheus.NewDesc(
			"api_request_count_total",
			"API请求总次数",
			nil,
			prometheus.Labels{"zone": zone},
		),
	}
}

// Describe 向Prometheus描述收集的指标
func (c *APIRequestCounter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.APIRequestDesc
}

// Collect 收集指标数据并发送到Prometheus
func (c *APIRequestCounter) Collect(ch chan<- prometheus.Metric) {
	count := atomic.LoadUint64(&c.requestCount)
	ch <- prometheus.MustNewConstMetric(
		c.APIRequestDesc,
		prometheus.CounterValue,
		float64(count),
	)
}

// IncrementRequestCount  增加API请求计数
func (c *APIRequestCounter) IncrementRequestCount() {
	atomic.AddUint64(&c.requestCount, 1)
}
