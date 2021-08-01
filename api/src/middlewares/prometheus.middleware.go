package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("use PrometheusMiddleware")
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(c.Request.URL.Path))
		totalRequests.WithLabelValues(c.Request.URL.Path).Inc()
		totalHTTPMethods.WithLabelValues(c.Request.Method).Inc()
		c.Next()
		timer.ObserveDuration()
	}
}
var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of incoming requests",
	},
	[]string{"path"},
)

var totalHTTPMethods = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_methods_total",
		Help: "Number of requests per HTTP method",
	},
	[]string{"method"},
)

var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests",
	},
	[]string{"path"},
)
func init() {
	log.Info("Init PrometheusMiddleware")
	prometheus.Register(totalRequests)
	prometheus.Register(totalHTTPMethods)
	prometheus.Register(httpDuration)
}