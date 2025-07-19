package handler

import (
	"komiko/config"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetCover(c *gin.Context) {
	c.File(filepath.Join(config.CoverDir, c.Param("file")))
}
