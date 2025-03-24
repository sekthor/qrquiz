package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var requestCounter metric.Int64Counter

func initCounter(meter metric.Meter) {
	var err error
	requestCounter, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		panic(err)
	}
}

func RequestCounter() gin.HandlerFunc {
	if requestCounter == nil {
		initCounter(otel.Meter("qrquiz"))
	}
	return func(c *gin.Context) {
		c.Next()
		requestCounter.Add(context.Background(), 1,
			metric.WithAttributes(
				attribute.String("http.request.method", c.Request.Method),
				attribute.String("http.url.scheme", c.Request.URL.Scheme),
				attribute.Int("http.response.status_code", c.Writer.Status()),
				attribute.String("http.route", c.Request.URL.Path),
			),
		)
	}
}
