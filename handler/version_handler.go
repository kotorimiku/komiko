package handler

import (
	"komiko/response"
	"komiko/version"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionHandler struct{}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

func (h *VersionHandler) GetVersion(c *gin.Context) {
	response.Success(c, version.GetShortVersion())
}

func (h *VersionHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": version.GetShortVersion(),
	})
}
