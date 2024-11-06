package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ridhotamma/yourkasa/product-service/config"
	"github.com/ridhotamma/yourkasa/product-service/routes"
	"github.com/ridhotamma/yourkasa/product-service/utils"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))

		c.Next()

		timer.ObserveDuration()
		status := string(rune(c.Writer.Status()))
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}

func main() {
	db := config.InitDB()
	r := gin.Default()

	r.Use(prometheusMiddleware())
	r.Use(utils.ValidationMiddleware())

	routes.SetupRoutes(r, db)

	if err := r.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
