package handler

import (
	"komiko/response"
	"komiko/service"
	"komiko/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type NovelHandler struct {
	NovelService *service.NovelService
}

func NewNovelHandler(services *service.Service) *NovelHandler {
	return &NovelHandler{
		NovelService: services.NovelService,
	}
}

func (h *NovelHandler) GetChapter(c *gin.Context) {
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
	data, err := h.NovelService.GetByPage(bookID, page)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	c.Data(http.StatusOK, http.DetectContentType(data), data)
}

func (h *NovelHandler) GetImage(c *gin.Context) {
	bookID, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	path := strings.TrimPrefix(c.Param("path"), "/")
	data, err := h.NovelService.GetByPath(bookID, path)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	c.Data(http.StatusOK, utils.ImageMiMeType(data), data)
}

func (h *NovelHandler) GetFile(c *gin.Context) {
	bookID, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	path := strings.TrimPrefix(c.Param("path"), "/")
	data, err := h.NovelService.GetByPath(bookID, path)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	t := http.DetectContentType(data)
	if strings.HasSuffix(path, ".xhtml") {
		t = "application/xhtml+xml"
	} else if strings.HasSuffix(path, ".css") {
		c.Header("Cache-Control", "public, max-age=86400")
		t = "text/css"
	} else if strings.HasSuffix(path, ".js") {
		c.Header("Cache-Control", "public, max-age=86400")
		t = "text/javascript"
	} else if strings.HasSuffix(path, ".ttf") {
		c.Header("Cache-Control", "public, max-age=86400")
	}

	c.Data(http.StatusOK, t, data)
}

func (h *NovelHandler) GetChapters(c *gin.Context) {
	bookID, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}

	chapters, err := h.NovelService.GetChapters(bookID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, chapters)
}
