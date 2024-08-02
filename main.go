package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-redis-poc/metrics_controller"
	"log"
	"net/http"
	"time"
)

var (
	rdb              *redis.Client
	ctx              = context.Background()
	keyPrefix        = "go-redis-demo:prefix:"
	metricsNamespace = flag.String("metric.namespace", "app", "Prometheus metrics namespace, as the prefix of metrics name")
)

func main() {
	flag.Parse()
	apiRequestCounter := metrics_controller.NewAPIRequestCounter(*metricsNamespace) // 创建一个新的Prometheus指标注册表
	registry := prometheus.NewRegistry()
	// 注册APIRequestCounter实例到Prometheus注册表
	registry.MustRegister(apiRequestCounter)

	clusterSlots := func(ctx context.Context) ([]redis.ClusterSlot, error) {
		slots := []redis.ClusterSlot{
			// First node with 1 master and 1 slave.
			{
				Start: 0,
				End:   5460,
				Nodes: []redis.ClusterNode{{
					Addr: "192.168.31.143:7001", // master
				}, {
					Addr: "192.168.31.143:7005", // 1st slave
				}},
			},
			// Second node with 1 master and 1 slave.
			{
				Start: 5461,
				End:   10922,
				Nodes: []redis.ClusterNode{{
					Addr: "192.168.31.143:7002", // master
				}, {
					Addr: "192.168.31.143:7006", // 1st slave
				}},
			},
			{
				Start: 10923,
				End:   16383,
				Nodes: []redis.ClusterNode{{
					Addr: "192.168.31.143:7003", // master
				}, {
					Addr: "192.168.31.143:7004", // 1st slave
				}},
			},
		}
		return slots, nil
	}

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots:  clusterSlots,
		RouteRandomly: true,
	})
	rdb.Ping(ctx)

	// ReloadState reloads cluster state. It calls ClusterSlots func
	// to get cluster slots information.
	rdb.ReloadState(ctx)

	// Initialize Gin engine
	r := gin.Default()
	r.GET("/", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
	            <head><title>A Prometheus Exporter</title></head>
	            <body>
	            <h1>A Prometheus Exporter</h1>
	            <p><a href='/metrics'>Metrics</a></p>
	            </body>
	            </html>`))
	}))

	// 模拟API请求的处理函数
	r.GET("/api", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		apiRequestCounter.IncrementRequestCount()
		// 模拟API处理时间
		// api handler的逻辑处理
		time.Sleep(1 * time.Second)

		w.Write([]byte("API请求处理成功"))
	}))

	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))
	// 设置根路径的处理函数，用于返回一个简单的HTML页面，包含指向指标页面的链接

	// Set key-value pair in Redis
	r.POST("/set/:key/:value", func(c *gin.Context) {
		key := keyPrefix + c.Param("key") // Prefix keys with 'myapp:'
		value := c.Param("value")
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			log.Printf("Error setting key-value pair in Redis: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Key-value pair set successfully in Redis: %s = %s\n", key, value)
		c.JSON(http.StatusOK, gin.H{"message": "Key-value pair set successfully"})
	})

	// Get value of key from Redis
	r.GET("/get/:key", func(c *gin.Context) {
		key := keyPrefix + c.Param("key") // Prefix keys with 'myapp:'
		value, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			log.Printf("Key not found in Redis: %s\n", key)
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		} else if err != nil {
			log.Printf("Error getting value from Redis: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Value retrieved from Redis: %s = %s\n", key, value)
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	// Start Gin server
	r.Run(":8080")
}
