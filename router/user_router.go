package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewUserHandler(service)

	userPublicGroup := publicGroup.Group("/user")
	{
		userPublicGroup.POST("/login", handler.Login)
		userPublicGroup.POST("/register", handler.Register)
		userPublicGroup.GET("/allow-register", handler.AllowRegister)
	}

	userAuthGroup := authGroup.Group("/user")
	{
		userAuthGroup.GET("", handler.GetAll)
		userAuthGroup.GET("/:id", handler.GetByID)
		userAuthGroup.POST("", handler.Create)
		userAuthGroup.POST("/update", handler.Update)
		userAuthGroup.POST("/:id/delete", handler.DeleteByID)

		userAuthGroup.GET("/current", handler.GetCurrentUser)
	}

}
