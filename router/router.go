package router

import (
	"komiko/middleware"
	"komiko/repo"
	"komiko/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	publicGroup := r.Group("/api")
	authGroup := r.Group("/api")
	repo := repo.NewRepo(db)
	authGroup.Use(middleware.AuthMiddleware(repo))
	service := service.NewService(repo)
	RegisterBookRoutes(authGroup, publicGroup, service)
	RegisterLibraryRoutes(authGroup, publicGroup, service)
	RegisterSeriesRoutes(authGroup, publicGroup, service)
	RegisterCoverRoutes(authGroup, publicGroup)
	RegisterUserRoutes(authGroup, publicGroup, service)
	RegisterProgressRoutes(authGroup, publicGroup, service)
	RegisterTaskRoutes(authGroup, publicGroup)
}
