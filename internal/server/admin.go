package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sekthor/qrquiz/internal/config"
)

func requiresAdmin(config *config.Config) gin.HandlerFunc {
	if config.Admin.Disabled {
		return func(*gin.Context) {}
	}
	return gin.BasicAuth(gin.Accounts{
		config.Admin.User: config.Admin.Password,
	})
}
