package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterSeriesRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewSeriesHandler(service)

	seriesGroup := authGroup.Group("/series")
	{
		seriesGroup.GET("", handler.GetList)
		seriesGroup.GET("/:id", handler.GetByID)
	}

}
