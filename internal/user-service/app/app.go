package app

// import (
// 	"hr-program/internal/user-service/app/router"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/prometheus/client_golang/prometheus"
// )

// func NewApp() *gin.Engine {
// 	app := gin.New()
// 	app.Use(PrometheusMiddleware())
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "*",
// 		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 		// ExposeHeaders:    []string{"Content-Length", "Content-Type"},
// 		// AllowCredentials: true,
// 	}))
// 	router.UserRouter(app)
// 	return app
// }

// var (
// 	httpRequests = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "user_http_requests_total",
// 			Help: "Total number of HTTP requests",
// 		},
// 		[]string{"route_group", "route", "method", "status"},
// 	)

// 	requestDuration = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "user_http_request_duration_seconds",
// 			Help:    "HTTP request duration in seconds",
// 			Buckets: prometheus.DefBuckets,
// 		},
// 		[]string{"route_group", "route", "method"},
// 	)
// )

// func Register() {
// 	prometheus.MustRegister(httpRequests)
// 	prometheus.MustRegister(requestDuration)
// }

// func PrometheusMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		route := c.Path()
// 		if route == "" {
// 			route = "unknown"
// 		}
// 		method := c.Method()

// 		err := c.Next()

// 		duration := time.Since(start).Seconds()
// 		status := strconv.Itoa(c.Response().StatusCode())

// 		httpRequests.WithLabelValues("user", route, method, status).Inc()
// 		requestDuration.WithLabelValues("user", route, method, status).Observe(duration)
// 		return err
// 	}
// }
