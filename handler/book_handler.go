package handler

import (
	"komiko/model"
	"komiko/response"
	"komiko/service"
	"komiko/utils"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	baseHandler[model.Book, *service.BookService]
}

func NewBookHandler(services *service.Service) *BookHandler {
	return &BookHandler{
		baseHandler: baseHandler[model.Book, *service.BookService]{
			service: services.BookService,
		},
	}
}

func (h *BookHandler) GetBySeriesID(c *gin.Context) {
	seriesId, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	books, err := h.service.GetBySeriesID(seriesId)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, books)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	seriesID := c.Query("series")

	books, err := h.service.GetBooks(seriesID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, books)
}
