package router

import (
	"komiko/handler"

	"github.com/gin-gonic/gin"
)

func RegisterVersionRoutes(publicGroup *gin.RouterGroup) {
	versionHandler := handler.NewVersionHandler()
	publicGroup.GET("/version", versionHandler.GetVersion)
	publicGroup.GET("/health", versionHandler.HealthCheck)
}
