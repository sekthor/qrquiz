package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// Process the request
		c.Next()

		statusCode := c.Writer.Status()

		// Log the request details
		entry := logrus.WithFields(logrus.Fields{
			"status":  statusCode,
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"proto":   c.Request.Proto,
			"latency": time.Since(start),
		}).WithContext(c.Request.Context())

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.URL.Path)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
