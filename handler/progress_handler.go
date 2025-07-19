package handler

import (
	"komiko/model"
	"komiko/response"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	baseHandler[model.Progress, *service.ProgressService]
}

func NewProgressHandler(services *service.Service) *ProgressHandler {
	return &ProgressHandler{
		baseHandler: baseHandler[model.Progress, *service.ProgressService]{
			service: services.ProgressService,
		},
	}
}

func (h *ProgressHandler) GetBookProgresses(c *gin.Context) {
	userID := c.GetUint("userID")
	seriesIdStr := c.Query("series")
	limit := c.Query("limit")
	offset := c.Query("offset")

	progress, err := h.service.GetBookProgresses(userID, seriesIdStr, limit, offset)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, progress)
}

func (h *ProgressHandler) GetSeriesProgresses(c *gin.Context) {
	userID := c.GetUint("userID")
	libraryIdStr := c.Query("library")
	limit := c.Query("limit")
	offset := c.Query("offset")

	progress, err := h.service.GetSeriesProgresses(userID, libraryIdStr, limit, offset)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, progress)
}

func (h *ProgressHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	var progress model.Progress
	if err := c.ShouldBindJSON(&progress); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	progress.UserID = userID
	if err := h.service.UpdateOrCreateByBookID(&progress); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, progress)
}
