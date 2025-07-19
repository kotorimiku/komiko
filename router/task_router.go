package router

import (
	"komiko/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {
	taskGroup := publicGroup.Group("/task")
	{
		taskGroup.GET("/:id", handler.GetTaskHandler)
		taskGroup.POST("/:id/stop", handler.StopTaskHandler)
		taskGroup.GET("/", handler.ListTasksHandler)
	}
}
