package handler

import (
	"komiko/response"
	"komiko/service"
	"komiko/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComicHandler struct {
	ComicService *service.ComicService
}

func NewComicHandler(services *service.Service) *ComicHandler {
	return &ComicHandler{
		ComicService: services.ComicService,
	}
}

func (h *ComicHandler) GetImage(c *gin.Context) {
	bookID, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	page, err := utils.ParamToUint(c.Param("page"))
	if err != nil {
		response.Error(c, "Invalid page")
		return
	}
	data, err := h.ComicService.GetByPage(bookID, page)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	c.Data(http.StatusOK, utils.ImageMiMeType(data), data)
}
