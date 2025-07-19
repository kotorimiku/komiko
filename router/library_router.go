package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterLibraryRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewLibraryHandler(service)

	libraryGroup := authGroup.Group("/library")
	{
		libraryGroup.GET("", handler.GetAll)
		libraryGroup.GET("/:id", handler.GetByID)
		libraryGroup.POST("", handler.Create)
		libraryGroup.POST("/update", handler.Update)
		libraryGroup.POST("/:id/delete", handler.DeleteByID)
		libraryGroup.POST("/:id/update-cover", handler.UpdateCover)
		libraryGroup.POST("/:id/scan-update", handler.ScanUpdate)
		libraryGroup.POST("/:id/scan-create", handler.ScanCreate)
	}

}
